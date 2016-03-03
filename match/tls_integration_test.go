package match

import (
	"bytes"
	"testing"
)

var testTLSPackets = map[string][]byte{
	"curl/7.44.0 OpenSSL/1.0.2e": []byte{0x16, 0x03, 0x01, 0x00, 0xb9, 0x01, 0x00, 0x00, 0xb5, 0x03, 0x03, 0x56, 0xd8, 0xaf, 0x3e, 0x0b,
		0x63, 0xb7, 0xd0, 0xcf, 0xfb, 0xfd, 0xd4, 0x09, 0x37, 0xe4, 0xf8, 0xa7, 0x85, 0x7a, 0xbb, 0xe0,
		0xc0, 0x32, 0xf8, 0xde, 0x1b, 0x6a, 0x2c, 0xe8, 0xdf, 0x4b, 0xae, 0x00, 0x00, 0x56, 0x00, 0xff,
		0xc0, 0x24, 0xc0, 0x23, 0xc0, 0x0a, 0xc0, 0x09, 0xc0, 0x08, 0xc0, 0x28, 0xc0, 0x27, 0xc0, 0x14,
		0xc0, 0x13, 0xc0, 0x12, 0xc0, 0x26, 0xc0, 0x25, 0xc0, 0x05, 0xc0, 0x04, 0xc0, 0x03, 0xc0, 0x2a,
		0xc0, 0x29, 0xc0, 0x0f, 0xc0, 0x0e, 0xc0, 0x0d, 0x00, 0x6b, 0x00, 0x67, 0x00, 0x39, 0x00, 0x33,
		0x00, 0x16, 0x00, 0x3d, 0x00, 0x3c, 0x00, 0x35, 0x00, 0x2f, 0x00, 0x0a, 0xc0, 0x07, 0xc0, 0x11,
		0xc0, 0x02, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x04, 0x00, 0xaf, 0x00, 0xae, 0x00, 0x8d, 0x00, 0x8c,
		0x00, 0x8a, 0x00, 0x8b, 0x01, 0x00, 0x00, 0x36, 0x00, 0x00, 0x00, 0x10, 0x00, 0x0e, 0x00, 0x00,
		0x0b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x6f, 0x72, 0x67, 0x00, 0x0a, 0x00, 0x08,
		0x00, 0x06, 0x00, 0x17, 0x00, 0x18, 0x00, 0x19, 0x00, 0x0b, 0x00, 0x02, 0x01, 0x00, 0x00, 0x0d,
		0x00, 0x0c, 0x00, 0x0a, 0x05, 0x01, 0x04, 0x01, 0x02, 0x01, 0x04, 0x03, 0x02, 0x03},
	"chrome/48.0.2564.109 (64-bit)": []byte{0x16, 0x03, 0x01, 0x00, 0xc6, 0x01, 0x00, 0x00, 0xc2, 0x03, 0x03, 0xa4, 0xaa, 0xb0, 0xf2, 0x8d,
		0x08, 0x18, 0xb2, 0xba, 0x3d, 0x44, 0x7b, 0x24, 0x4f, 0xad, 0x0f, 0xaa, 0xe4, 0xe7, 0x67, 0x9a,
		0x61, 0x78, 0xc4, 0xed, 0x37, 0xd6, 0xf1, 0xd6, 0xfb, 0x43, 0xdb, 0x00, 0x00, 0x1e, 0xc0, 0x2b,
		0xc0, 0x2f, 0x00, 0x9e, 0xcc, 0x14, 0xcc, 0x13, 0xc0, 0x0a, 0xc0, 0x14, 0x00, 0x39, 0xc0, 0x09,
		0xc0, 0x13, 0x00, 0x33, 0x00, 0x9c, 0x00, 0x35, 0x00, 0x2f, 0x00, 0x0a, 0x01, 0x00, 0x00, 0x7b,
		0xff, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x0e, 0x00, 0x00, 0x0b, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x6f, 0x72, 0x67, 0x00, 0x17, 0x00, 0x00, 0x00, 0x23, 0x00,
		0x00, 0x00, 0x0d, 0x00, 0x16, 0x00, 0x14, 0x06, 0x01, 0x06, 0x03, 0x05, 0x01, 0x05, 0x03, 0x04,
		0x01, 0x04, 0x03, 0x03, 0x01, 0x03, 0x03, 0x02, 0x01, 0x02, 0x03, 0x00, 0x05, 0x00, 0x05, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x33, 0x74, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 0x10, 0x00, 0x17,
		0x00, 0x15, 0x02, 0x68, 0x32, 0x08, 0x73, 0x70, 0x64, 0x79, 0x2f, 0x33, 0x2e, 0x31, 0x08, 0x68,
		0x74, 0x74, 0x70, 0x2f, 0x31, 0x2e, 0x31, 0x75, 0x50, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x02, 0x01,
		0x00, 0x00, 0x0a, 0x00, 0x06, 0x00, 0x04, 0x00, 0x17, 0x00, 0x18},
	"safari/9.0.1 (10601.2.7.2)": []byte{0x16, 0x03, 0x01, 0x00, 0xb1, 0x01, 0x00, 0x00, 0xad, 0x03, 0x03, 0x56, 0xd8, 0xac, 0x9c, 0x0f,
		0x92, 0x85, 0xe7, 0xab, 0x32, 0x0e, 0xbe, 0xc4, 0x7e, 0x4e, 0xc7, 0x81, 0xcd, 0x79, 0x74, 0x56,
		0x57, 0x2c, 0x4a, 0x30, 0xb9, 0x96, 0x7f, 0xd8, 0xe4, 0x4d, 0xeb, 0x00, 0x00, 0x4a, 0x00, 0xff,
		0xc0, 0x24, 0xc0, 0x23, 0xc0, 0x0a, 0xc0, 0x09, 0xc0, 0x08, 0xc0, 0x28, 0xc0, 0x27, 0xc0, 0x14,
		0xc0, 0x13, 0xc0, 0x12, 0xc0, 0x26, 0xc0, 0x25, 0xc0, 0x05, 0xc0, 0x04, 0xc0, 0x03, 0xc0, 0x2a,
		0xc0, 0x29, 0xc0, 0x0f, 0xc0, 0x0e, 0xc0, 0x0d, 0x00, 0x6b, 0x00, 0x67, 0x00, 0x39, 0x00, 0x33,
		0x00, 0x16, 0x00, 0x3d, 0x00, 0x3c, 0x00, 0x35, 0x00, 0x2f, 0x00, 0x0a, 0xc0, 0x07, 0xc0, 0x11,
		0xc0, 0x02, 0xc0, 0x0c, 0x00, 0x05, 0x00, 0x04, 0x01, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x00, 0x10,
		0x00, 0x0e, 0x00, 0x00, 0x0b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x6f, 0x72, 0x67,
		0x00, 0x0a, 0x00, 0x08, 0x00, 0x06, 0x00, 0x17, 0x00, 0x18, 0x00, 0x19, 0x00, 0x0b, 0x00, 0x02,
		0x01, 0x00, 0x00, 0x0d, 0x00, 0x0c, 0x00, 0x0a, 0x05, 0x01, 0x04, 0x01, 0x02, 0x01, 0x04, 0x03,
		0x02, 0x03, 0x33, 0x74, 0x00, 0x00},
}

var testTLSPacketsHostname = []byte("example.org")

func TestTLSIntegration(t *testing.T) {
	for client, packet := range testTLSPackets {
		hostname, err := TLS(packet)
		if err != nil {
			t.Logf("client %s: error: %s", client, err)
		} else if !bytes.Equal(hostname, testTLSPacketsHostname) {
			t.Logf("client %s: error: got '%s' expected '%s'", client, hostname, testTLSPacketsHostname)
		}
	}
}
