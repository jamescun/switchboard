package server

import (
	"bytes"
	"io"
	"testing"

	"github.com/jamescun/switchboard/match"

	"github.com/stretchr/testify/assert"
)

func TestInitBufPool(t *testing.T) {
	s := &Server{}
	s.initBufPool()
	assert.NotNil(t, s.buf.New)
}

func TestServerMatch(t *testing.T) {
	s := &Server{Match: match.Http}

	_, _, err := s.match(bytes.NewReader([]byte("GET / HTTP/1.1\nHost: example.org\n\n")), make([]byte, 14))
	assert.Equal(t, io.ErrShortBuffer, err)

	hostname, n, err := s.match(bytes.NewReader([]byte("GET / HTTP/1.1\nHost: example.org\n\n")), make([]byte, 64))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("example.org"), hostname)
		assert.Equal(t, 34, n)
	}
}

func TestServerProxy(t *testing.T) {
	s := &Server{}

	rx := make([]byte, 64)
	tx := make([]byte, 64)

	n := copy(rx, []byte("GET / HTTP/1.1\nHost: example.org\n\n"))

	local := &conn{rx: bytes.NewBuffer(nil), tx: bytes.NewBuffer(nil)}
	remote := &conn{rx: bytes.NewBuffer(nil), tx: bytes.NewBuffer([]byte("HTTP/1.0 404 Not Found\n"))}

	err := s.proxy(remote, local, rx, tx, n)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("GET / HTTP/1.1\nHost: example.org\n\n"), remote.rx.Bytes())
		assert.Equal(t, []byte("HTTP/1.0 404 Not Found\n"), local.rx.Bytes())
	}
}
