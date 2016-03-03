package match

import (
	"bytes"
	"errors"
	"unicode"
)

var (
	ErrNone       = errors.New("no match found")
	ErrShortBytes = errors.New("short bytes")
)

// Match is a function which implements the extraction of a hostname from a received packet.
// the resulting bytes MAY be a slice of the original packet and MUST NOT be modified.
// if the packet contains no or a partial match, the function MUST return ErrNone. the caller
// can then decide whether to continue, read more data or terminate the connection.
// all other errors must terminate the connection.
type Match func(packet []byte) (hostname []byte, err error)

// return index of first matching bytes
func indexFirst(b []byte, v ...[]byte) int {
	for i := 0; i < len(v); i++ {
		j := bytes.Index(b, v[i])
		if j > -1 {
			return j
		}
	}

	return -1
}

// return slice of bytes skipped forward of any whitespace
func skip(b []byte) (o []byte) {
	if len(b) == 0 {
		return
	}

	i := bytes.IndexFunc(b, notSpace)
	if i > -1 {
		o = b[i:]
	}

	return
}

// return true if rune is not whitespace character
func notSpace(r rune) bool {
	return !unicode.IsSpace(r)
}
