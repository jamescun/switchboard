package backend

// Static implements the Backend interface only to return a static set of pre-configured
// upstream addresses.
type Static struct {
	Upstreams []string
}

func (s Static) Start() error                        { return nil }
func (s Static) Stop() error                         { return nil }
func (s Static) Upstream(_ []byte) ([]string, error) { return s.Upstreams, nil }
