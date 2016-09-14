package mycipher

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

func NewHash(key []byte) hash.Hash {
	return hmac.New(sha256.New, key)
}
