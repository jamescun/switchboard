package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	var s []byte
	var err error

	_, err = Http([]byte{})
	assert.Equal(t, ErrNone, err)

	_, err = Http([]byte("GET / HTTP/1.0\n\n"))
	assert.Equal(t, ErrNone, err)

	_, err = Http([]byte("GET / HTTP/1.1\nHost:"))
	assert.Equal(t, ErrNone, err)

	_, err = Http([]byte("GET / HTTP/1.1\nHost: "))
	assert.Equal(t, ErrNone, err)

	_, err = Http([]byte("GET / HTTP/1.1\nHost: example.org"))
	assert.Equal(t, ErrNone, err)

	s, err = Http([]byte("GET / HTTP/1.1\nHost: example.org\nConnection: close\n\n"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("example.org"), s)
	}

	s, err = Http([]byte("GET / HTTP/1.1\nhost:example.org\r\n\n"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("example.org"), s)
	}

	s, err = Http([]byte("GET / HTTP/1.1\nHOST:  example.org\n\n"))
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("example.org"), s)
	}
}
