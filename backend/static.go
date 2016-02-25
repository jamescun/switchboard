package backend

// Static implements the Backend interface only to return a static set of pre-configured
// upstream addresses.
type Static struct {
	Upstreams []string

	u []Upstream
}

func (s *Static) Start() error {
	for i := 0; i < len(s.Upstreams); i++ {
		s.u = append(s.u, upstream(s.Upstreams[i]))
	}

	return nil
}

func (s *Static) Stop() error { return nil }

func (s *Static) Upstream(_ []byte) ([]Upstream, error) {
	if len(s.u) < 1 {
		return nil, ErrNone
	}

	return s.u, nil
}
