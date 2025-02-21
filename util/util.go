package util

import "math/rand"

func RandString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(26) + 65)
	}
	return string(b)
}
