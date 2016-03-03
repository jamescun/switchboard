package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTLSPacket = []byte{
	// TLS record
	0x16,       // Content Type: Handshake
	0x03, 0x01, // Version: TLS 1.0
	0x00, 0x6c, // Length (use for bounds checking)
	// Handshake
	0x01,             // Handshake Type: Client Hello
	0x00, 0x00, 0x68, // Length (use for bounds checking)
	0x03, 0x03, // Version: TLS 1.2
	// Random (32 bytes fixed length)
	0xb6, 0xb2, 0x6a, 0xfb, 0x55, 0x5e, 0x03, 0xd5,
	0x65, 0xa3, 0x6a, 0xf0, 0x5e, 0xa5, 0x43, 0x02,
	0x93, 0xb9, 0x59, 0xa7, 0x54, 0xc3, 0xdd, 0x78,
	0x57, 0x58, 0x34, 0xc5, 0x82, 0xfd, 0x53, 0xd1,
	0x00,       // Session ID Length (skip past this much)
	0x00, 0x04, // Cipher Suites Length (skip past this much)
	0x00, 0x01, // NULL-MD5
	0x00, 0xff, // RENEGOTIATION INFO SCSV
	0x01,       // Compression Methods Length (skip past this much)
	0x00,       // NULL
	0x00, 0x3b, // Extensions Length (use for bounds checking)
	// Extension
	0x00, 0x00, // Extension Type: Server Name (check extension type)
	0x00, 0x0e, // Length (use for bounds checking)
	0x00, 0x0c, // Server Name Indication Length
	0x00,       // Server Name Type: host_name (check server name type)
	0x00, 0x09, // Length (length of your data)
	// "localhost" (data your after)
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74,
	// Extension
	0x00, 0x0d, // Extension Type: Signature Algorithms (check extension type)
	0x00, 0x20, // Length (skip past since this is the wrong extension)
	// Data
	0x00, 0x1e, 0x06, 0x01, 0x06, 0x02, 0x06, 0x03,
	0x05, 0x01, 0x05, 0x02, 0x05, 0x03, 0x04, 0x01,
	0x04, 0x02, 0x04, 0x03, 0x03, 0x01, 0x03, 0x02,
	0x03, 0x03, 0x02, 0x01, 0x02, 0x02, 0x02, 0x03,
	// Extension
	0x00, 0x0f, // Extension Type: Heart Beat (check extension type)
	0x00, 0x01, // Length (skip past since this is the wrong extension)
	0x01, // Mode: Peer allows to send requests
}

func TestTLS(t *testing.T) {
	var o []byte
	var err error

	_, err = TLS(testTLSPacket[:43])
	assert.Equal(t, ErrShortBytes, err)

	_, err = TLS(testTLSPacket[:45])
	assert.Equal(t, ErrShortBytes, err)

	_, err = TLS(testTLSPacket[:51])
	assert.Equal(t, ErrShortBytes, err)

	_, err = TLS(testTLSPacket[:53])
	assert.Equal(t, ErrShortBytes, err)

	_, err = TLS(testTLSPacket[:54])
	assert.Equal(t, ErrShortBytes, err)

	o, err = TLS(testTLSPacket)
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("localhost"), o)
	}
}

func BenchmarkTLS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TLS(testTLSPacket)
	}
}

func TestReadSNIHost(t *testing.T) {
	var s []byte
	var err error

	_, err = readSNIHost([]byte{0x00, 0x00})
	assert.Equal(t, ErrShortBytes, err)

	_, err = readSNIHost([]byte{0x00, 0x00, 0x00, 0x00})
	assert.Equal(t, ErrShortBytes, err)

	_, err = readSNIHost([]byte{0x00, 0x00, 0x00, 0x00, 0x09})
	assert.Equal(t, ErrShortBytes, err)

	s, err = readSNIHost([]byte{0x00, 0x0c, 0x00, 0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74}, s)
	}
}

func BenchmarkReadSNIHost(b *testing.B) {
	p := []byte{0x00, 0x0c, 0x00, 0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readSNIHost(p)
	}
}

func TestReadUint16(t *testing.T) {
	var n int
	var o []byte
	var err error

	_, _, err = readUint16([]byte{0x01})
	assert.Equal(t, ErrShortBytes, err)

	n, o, err = readUint16([]byte{0x00, 0x09, 0xFF})
	if assert.NoError(t, err) {
		assert.Equal(t, 9, n)
		assert.Equal(t, []byte{0xFF}, o)
	}
}

func BenchmarkReadUint16(b *testing.B) {
	p := []byte{0x00, 0x09, 0x0FF}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readUint16(p)
	}
}

func TestSkipUint8(t *testing.T) {
	var o []byte
	var err error

	_, err = skipUint8([]byte{})
	assert.Equal(t, ErrShortBytes, err)

	_, err = skipUint8([]byte{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73})
	assert.Equal(t, ErrShortBytes, err)

	o, err = skipUint8([]byte{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte{}, o)
	}

	o, err = skipUint8([]byte{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0xFF})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0xFF}, o)
	}
}

func BenchmarkSkipUint8(b *testing.B) {
	p := []byte{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0xFF}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skipUint8(p)
	}
}

func TestSkipUint16(t *testing.T) {
	var o []byte
	var err error

	_, err = skipUint16([]byte{0x00})
	assert.Equal(t, ErrShortBytes, err)

	o, err = skipUint16([]byte{0x00, 0x00})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte(nil), o)
	}

	_, err = skipUint16([]byte{0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73})
	assert.Equal(t, ErrShortBytes, err)

	o, err = skipUint16([]byte{0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte{}, o)
	}

	o, err = skipUint16([]byte{0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0xFF})
	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0xFF}, o)
	}
}

func BenchmarkSkipUint16(b *testing.B) {
	p := []byte{0x00, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0xFF}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skipUint16(p)
	}
}
