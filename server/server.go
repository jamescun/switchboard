package server

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/jamescun/switchboard/backend"
	"github.com/jamescun/switchboard/match"
)

type Server struct {
	Match      match.Match
	Backend    backend.Backend
	BufferSize int
}

// Serve accepts incomming connections on the Listener l, creating a
// new service goroutine for each. the service goroutine will read
// packets from the client until a hostname match is found, initial
// buffer is full or client times out.
// Serve always returns a non-nil error.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	var tmpDelay time.Duration // backoff timeout for Accept() failures
	for {
		conn, err := l.Accept()
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			// if the error is temporary, backoff exponentially until error
			// is resolved; useful in cases like running out of file descriptors.
			if tmpDelay == 0 {
				tmpDelay = 5 * time.Millisecond
			} else {
				tmpDelay *= 2
			}

			if max := 1 * time.Second; tmpDelay > max {
				tmpDelay = max
			}

			time.Sleep(tmpDelay)
			continue
		} else if err != nil {
			return err
		}
		tmpDelay = 0

		go s.serve(conn)
	}

	return nil
}

// TODO: smell: break down, optimise and test
func (s *Server) serve(local net.Conn) {
	defer local.Close()

	// TODO: optimisation: use sync.Pool to reuse client buffers
	rx := make([]byte, s.BufferSize)
	tx := make([]byte, s.BufferSize)

	hostname, n, err := s.match(local, rx)
	if err == match.ErrNone || err == io.ErrShortBuffer {
		log.Println("warning: match: no match found")
		return
	}
	log.Printf("info: match: '%s'\n", hostname)

	upstream, err := s.Backend.Upstream(hostname)
	if err == backend.ErrNone || len(upstream) < 1 {
		log.Println("warning: backend: upstream: no match")
		return
	} else if err != nil {
		log.Println("error: backend: upstream:", err)
		return
	}

	remote, err := net.Dial("tcp", upstream[0].Addr())
	if err != nil {
		log.Println("error: remote:", err)
		return
	}
	defer remote.Close()

	err = s.proxy(remote, local, rx, tx, n)
	if err != nil {
		log.Println("error: proxy:", err)
	}
}

// return the hostname matched and the number of bytes read from reader into rx buffer.
// hostname is a byte slice of the initial packet.
func (s *Server) match(r io.Reader, rx []byte) (hostname []byte, n int, err error) {
	for {
		if n >= len(rx) {
			// buffer was not big enough to match hostname
			err = io.ErrShortBuffer
			return
		}

		j, rerr := r.Read(rx[n:])
		if rerr != nil {
			err = rerr
			return
		}
		n += j

		// use math.Match function to extract hostname from initial packet(s)
		hostname, rerr = s.Match(rx[:n])
		if rerr == match.ErrNone {
			// match not found, fill buffer some more
			continue
		} else if rerr != nil {
			err = rerr
			return
		}

		break
	}

	return
}

// bi-directional proxy between local and remote using rx/tx buffers, retransmitting the
// initial packet to the remote first.
func (s *Server) proxy(remote, local io.ReadWriter, rx, tx []byte, n int) error {
	// send initial packet to upstream
	_, err := remote.Write(rx[:n])
	if err != nil {
		return err
	}

	return ProxyBuffer(remote, local, rx, tx)
}
