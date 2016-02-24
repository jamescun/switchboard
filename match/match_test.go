package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkHttp(b *testing.B) {
	p := []byte("GET / HTTP/1.1\nHost: example.org\nConnection: close\n\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Http(p)
	}
}

func TestIndexFirst(t *testing.T) {
	assert.Equal(t, indexFirst([]byte("foo"), []byte("o")), 1)
	assert.Equal(t, indexFirst([]byte("bar"), []byte("o"), []byte("r")), 2)
	assert.Equal(t, indexFirst([]byte("baz"), []byte("j")), -1)
}

func BenchmarkIndexFirst(b *testing.B) {
	p := []byte("GET / HTTP/1.1\nHost: example.org\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indexFirst(p, httpHdrHostRFC, httpHdrHostLower, httpHdrHostUpper)
	}
}

func TestSkip(t *testing.T) {
	var b []byte

	b = skip([]byte{})
	assert.Len(t, b, 0)

	b = skip([]byte("  "))
	assert.Len(t, b, 0)

	b = skip([]byte("hello"))
	assert.Equal(t, []byte("hello"), b)

	b = skip([]byte(" hello"))
	assert.Equal(t, []byte("hello"), b)

	b = skip([]byte("\t \thello"))
	assert.Equal(t, []byte("hello"), b)
}

func BenchmarkSkip(b *testing.B) {
	p := []byte(" example.org")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skip(p)
	}
}

func TestNotSpace(t *testing.T) {
	assert.True(t, notSpace('j'))
	assert.False(t, notSpace(' '))
}
