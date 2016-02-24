package server

import (
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

	var o int
	var hn []byte
	for {
		n, err := local.Read(rx[o:])
		if err != nil {
			log.Println("error: read:", err)
			return
		}
		o += n
		if o >= len(rx) {
			log.Println("error: read: buffer full")
			return
		}

		hn, err = s.Match(rx)
		if err == match.ErrNone {
			log.Println("warning: match: no match, retrying")
			continue
		} else if err != nil {
			log.Println("error: match:", err)
			return
		}

		log.Printf("info: match: '%s'\n", hn)
		break
	}

	upstream, err := s.Backend.Upstream(hn)
	if err == backend.ErrNone || len(upstream) < 1 {
		log.Println("warning: backend: upstream: no match")
		return
	} else if err != nil {
		log.Println("error: backend: upstream:", err)
		return
	}

	remote, err := net.Dial("tcp", upstream[0])
	if err != nil {
		log.Println("error: remote:", err)
		return
	}
	defer remote.Close()

	_, err = remote.Write(rx[:o])
	if err != nil {
		log.Println("error: write:", err)
		return
	}

	err = ProxyBuffer(remote, local, rx, tx)
	if err != nil {
		log.Println("error: proxy:", err)
		return
	}

	log.Println("info: proxy: complete")
}
