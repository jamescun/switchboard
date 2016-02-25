package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpstreamAddr(t *testing.T) {
	u := upstream("192.168.0.1:8080")
	assert.Equal(t, "192.168.0.1:8080", u.Addr())
}
