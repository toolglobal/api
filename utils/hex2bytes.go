package utils

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func HexToBytes(str string) []byte {
	str = strings.TrimPrefix(str, "0x")
	b, _ := hex.DecodeString(str)
	return b
}

func B64ToPubKey(str string) (h [32]byte) {
	b, _ := base64.StdEncoding.DecodeString(str)
	if len(b) > len(h) {
		b = b[len(b)-32:]
	}
	copy(h[32-len(b):], b)
	return h
}

func B64ToPrivKey(str string) (h [64]byte) {
	b, _ := base64.StdEncoding.DecodeString(str)
	if len(b) > len(h) {
		b = b[len(b)-64:]
	}
	copy(h[64-len(b):], b)
	return h
}
