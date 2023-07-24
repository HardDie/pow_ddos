package utils

import (
	"crypto/rand"
	"fmt"
	"math/bits"
)

func RandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error generate random buffer: %w", err)
	}
	if n != size {
		return nil, fmt.Errorf("error read random bytes, got=%d; want=%d", n, size)
	}
	return buf, nil
}

func CountLeadingZeros(data []byte) int {
	count := 0
	for _, d := range data {
		if d == 0 {
			count += 8
			continue
		}
		count += bits.LeadingZeros8(d)
		break
	}
	return count
}
