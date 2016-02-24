package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticUpstream(t *testing.T) {
	u := []string{"10.0.0.1:1993", "10.0.0.2:1993"}

	s := Static{Upstreams: u}
	n, err := s.Upstream([]byte("example.org"))
	if assert.NoError(t, err) {
		assert.Equal(t, u, n)
	}
}
