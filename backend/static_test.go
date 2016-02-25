package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticStart(t *testing.T) {
	s := Static{Upstreams: []string{"192.168.0.1:8080", "192.168.0.2:8080"}}
	err := s.Start()
	if assert.NoError(t, err) {
		assert.Equal(t, []Upstream{upstream("192.168.0.1:8080"), upstream("192.168.0.2:8080")}, s.u)
	}
}

func TestStaticStop(t *testing.T) {
	s := Static{}
	err := s.Stop()
	assert.NoError(t, err)
}

func TestStaticUpstream(t *testing.T) {
	s := Static{}
	s.Start()
	_, err := s.Upstream([]byte("example.org"))
	assert.Equal(t, ErrNone, err)

	s = Static{Upstreams: []string{"192.168.0.1:8080", "192.168.0.2:8080"}}
	s.Start()
	n, err := s.Upstream([]byte("example.org"))
	if assert.NoError(t, err) {
		assert.Equal(t, []Upstream{upstream("192.168.0.1:8080"), upstream("192.168.0.2:8080")}, n)
	}
}
