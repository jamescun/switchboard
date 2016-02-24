package server

import (
	"io"
)

// proxyBuffer copies from src to dst and vice versa until either EOF is reached or an error occurs,
// using supplied byte slices as buffers (similar to io.CopyBuffer). it returns the first error
// encountered while copying, if any.
func ProxyBuffer(dst, src io.ReadWriter, rx, tx []byte) error {
	errch := make(chan error, 2)

	go proxyCopy(dst, src, rx, errch)
	go proxyCopy(src, dst, tx, errch)

	err := <-errch
	if err != nil {
		return err
	}

	return nil
}

func proxyCopy(dst, src io.ReadWriter, buf []byte, errch chan error) {
	_, err := io.CopyBuffer(dst, src, buf)
	errch <- err
}
