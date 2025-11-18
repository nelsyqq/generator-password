package security

import (
	"crypto/rand"
	"math/big"
)

func GetRandomIndex(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}
