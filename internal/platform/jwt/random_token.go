package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

type RandomTokenGenerator struct {
	size int
}

func NewRandomTokenGenerator(size int) *RandomTokenGenerator {
	if size <= 0 {
		size = 32
	}

	return &RandomTokenGenerator{size: size}
}

func (g *RandomTokenGenerator) Generate() (string, error) {
	buffer := make([]byte, g.size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}
