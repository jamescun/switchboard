package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type conn struct {
	rx, tx *bytes.Buffer
}

func (c *conn) Read(b []byte) (int, error)  { return c.tx.Read(b) }
func (c *conn) Write(b []byte) (int, error) { return c.rx.Write(b) }

func TestProxyBuffer(t *testing.T) {
	rx := make([]byte, 64)
	tx := make([]byte, 64)

	local := &conn{rx: bytes.NewBuffer(nil), tx: bytes.NewBuffer([]byte("GET / HTTP/1.1\nHost: example.org\n\n"))}
	remote := &conn{rx: bytes.NewBuffer(nil), tx: bytes.NewBuffer([]byte("HTTP/1.0 404 Not Found\n"))}

	err := ProxyBuffer(remote, local, rx, tx)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("GET / HTTP/1.1\nHost: example.org\n\n"), remote.rx.Bytes())
		assert.Equal(t, []byte("HTTP/1.0 404 Not Found\n"), local.rx.Bytes())
	}
}
