package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
)

func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

type Address struct {
	ethcmn.Address
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.Hex()), nil
}

func ValidAddress(s string) bool {
	return ethcmn.IsHexAddress(s)
}
