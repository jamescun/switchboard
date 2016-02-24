package backend

import (
	"errors"
)

var (
	ErrNone = errors.New("no backend for host")
)

// Backend is an interface which implements the querying of a, potentially external, datastore
// for the upstream addresses which can respond to the incomming request. a Backend should not
// implement any caching or load balancing as these are configured at a lower level.
type Backend interface {
	// perform any pre-query required steps, such as launching goroutines or registering
	// with an external service.
	Start() error

	// perform any post-operation teardown, such as de-registering from an external service.
	Stop() error

	// returns host:port pairs representing upstream servers which can handle the incomming
	// request. if no upstream is available, it MUST return ErrNone.
	// all other errors must terminate the connection.
	Upstream(hostname []byte) (hosts []string, err error)
}
