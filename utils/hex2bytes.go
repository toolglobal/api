package utils

import (
	"encoding/base64"
	"encoding/hex"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
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

func RLPHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func ToWei(v *big.Int) *big.Int {
	return new(big.Int).Mul(v, big.NewInt(1e10))
}

func ToEther(v *big.Int) *big.Int {
	return new(big.Int).Div(v, big.NewInt(1e10))
}
