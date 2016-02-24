package match

import (
	"bytes"
	"unicode"
)

// supported styles of http header
var (
	httpHdrLn        = []byte("\n")
	httpHdrHostRFC   = []byte("\nHost:")
	httpHdrHostLower = []byte("\nhost:")
	httpHdrHostUpper = []byte("\nHOST:")
)

// Http implements a Match function for plaintext HTTP1.1+ requests. it respects the headers
// 'Host', 'host' and 'HOST'.
func Http(b []byte) (s []byte, err error) {
	if len(b) == 0 {
		err = ErrNone
		return
	}

	hb := indexFirst(b, httpHdrHostRFC, httpHdrHostLower, httpHdrHostUpper)
	if hb < 0 {
		err = ErrNone
		return
	}

	s = skip(b[hb+6:])
	if len(s) == 0 {
		err = ErrNone
		return
	}

	he := bytes.Index(s, httpHdrLn)
	if he < 0 {
		err = ErrNone
		return
	}
	s = bytes.TrimRightFunc(s[:he], unicode.IsSpace)

	return
}
