package backend

import (
	"errors"
)

var (
	ErrNone = errors.New("no backend for host")
)

// Backend is an interface which implements the querying of a, potentially external, datastore
// for the upstream addresses which can respond to the incomming request. a service Backend should
// not implement any caching or load balancing as these are configured at a higher level.
type Backend interface {
	// Start is called exactly once before any beginning the server and querying for upstreams.
	// this can be used to launch additional servers (i.e. webhook based event buses) or to
	// register the backends presence with an external service.
	Start() error

	// Stop will be called when the server is gracefully shutting down. like Start it must also
	// close down any additional servers or de-register with external services.
	Stop() error

	// Upstream takes a byte slice from the initial packet representing the requested hostname
	// and returns an array of Upstream interfaces of all services that can fulfil the request.
	// if no upstream is available, it MUST return ErrNone.
	Upstream(hostname []byte) ([]Upstream, error)
}

// Upstream is an interface which represents a potential upstream service, returned by service
// discovery, to be load balanced across.
type Upstream interface {
	// return the host:port pair for upstream, to be used with net.Dial.
	Addr() string
}

// simple wrapper for host:port strings to implement the Upstream interface.
type upstream string

func (u upstream) Addr() string {
	return string(u)
}
