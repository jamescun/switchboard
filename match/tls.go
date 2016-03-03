package match

import (
	"encoding/binary"
)

func TLS(pkt []byte) (hostname []byte, err error) {
	var n, l int
	if len(pkt) < 44 {
		err = ErrShortBytes
		return
	}
	pkt = pkt[43:]

	// session id
	pkt, err = skipUint8(pkt)
	if err != nil {
		return
	}

	// cipher suites
	pkt, err = skipUint16(pkt)
	if err != nil {
		return
	}

	// compression
	pkt, err = skipUint8(pkt)
	if err != nil {
		return
	}

	// extensions
	n, pkt, err = readUint16(pkt)
	if err != nil {
		return
	} else if n == 0 {
		return
	} else if len(pkt) < n {
		err = ErrShortBytes
		return
	}

	for {
		if len(pkt) < 4 {
			err = ErrShortBytes
			return
		}
		n, pkt, err = readUint16(pkt) // extension type
		if err != nil {
			return
		}
		l, pkt, err = readUint16(pkt) // extension length
		if err != nil {
			return
		} else if len(pkt) < l {
			err = ErrShortBytes
			return
		}

		if n == 0 {
			hostname, err = readSNIHost(pkt[:l])
			return
		} else {
			pkt = pkt[l:]
		}
	}

	return
}

func readSNIHost(ext []byte) (s []byte, err error) {
	// skip sni length and host type
	if len(ext) < 3 {
		err = ErrShortBytes
		return
	}
	ext = ext[3:]

	var n int
	n, ext, err = readUint16(ext)
	if err != nil {
		return
	}
	if len(ext) < n {
		err = ErrShortBytes
		return
	}

	s = ext[:n]
	return
}

func readUint16(b []byte) (n int, o []byte, err error) {
	if len(b) < 2 {
		err = ErrShortBytes
		return
	}

	n = int(binary.BigEndian.Uint16(b[:2]))
	o = b[2:]
	return
}

func skipUint8(b []byte) (o []byte, err error) {
	if len(b) < 1 {
		err = ErrShortBytes
		return
	}

	l := int(b[0]) + 1
	if len(b) < l {
		err = ErrShortBytes
		return
	}

	o = b[l:]
	return
}

func skipUint16(b []byte) (o []byte, err error) {
	var n int
	n, b, err = readUint16(b)
	if err != nil {
		return
	} else if n == 0 {
		return
	} else if len(b) < n {
		err = ErrShortBytes
		return
	}

	o = b[n:]
	return
}
