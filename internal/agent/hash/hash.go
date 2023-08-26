package hash

import (
	"crypto/sha256"
)

func MakeHash(inStr []byte) []byte {
	h := sha256.New()
	h.Write(inStr)
	return h.Sum(nil)
}
