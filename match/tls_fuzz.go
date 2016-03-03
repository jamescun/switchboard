// +build gofuzz

package match

func Fuzz(data []byte) int {
	hostname, err := TLS(data)
	if err != nil || len(hostname) < 1 {
		return 0
	}

	return 1
}
