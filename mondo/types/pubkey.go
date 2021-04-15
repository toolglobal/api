package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	//hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(PublicKey{})
)

const PubKeyLength = 33

type PublicKey [PubKeyLength]byte

func (pk PublicKey) Equals(pk1 PublicKey) bool {
	return pk == pk1
}

func (pk PublicKey) VerifyBytes(msg, signature []byte) bool {
	if len(signature) != 65 {
		return false
	}

	pubkey, err := crypto.SigToPub(msg, signature)
	if err != nil {
		return false
	}
	b := crypto.CompressPubkey(pubkey)
	if ok := bytes.Equal(b, pk.Bytes()); !ok {
		return false
	}

	return true
}

// unsafe 兼容老版本
func (pk PublicKey) ToAddress() Address {
	if pk.IsAddress() {
		var addr Address
		copy(addr.Address[:], pk[13:])
		return addr
	}

	pubKey, err := crypto.DecompressPubkey(pk[:])
	if err != nil {
		return Address{}
	}
	addr := Address{
		crypto.PubkeyToAddress(*pubKey),
	}
	return addr
}

// unsafe 前13个字节都是0
func (a PublicKey) IsAddress() bool {
	for i := 0; i < 13; i++ {
		if a[i] != 0 {
			return false
		}
	}
	return true
}

func (a PublicKey) Bytes() []byte { return a[:] }

func (a PublicKey) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

func (a PublicKey) Hex() string {
	return hexutil.Encode(a[:])
}

// String implements fmt.Stringer.
func (a PublicKey) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a PublicKey) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *PublicKey) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-PubKeyLength:]
	}
	copy(a[PubKeyLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a PublicKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *PublicKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *PublicKey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

func HexToPubkey(s string) (PublicKey, error) {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, err
	}
	if len(b) != 33 {
		return PublicKey{}, errors.New("wrong pk len")
	}

	if _, err := crypto.DecompressPubkey(b[:]); err != nil {
		return PublicKey{}, err
	}

	var pubKey PublicKey
	pubKey.SetBytes(b)
	return pubKey, nil
}

func BytesToPubkey(b []byte) (PublicKey, error) {
	if _, err := crypto.DecompressPubkey(b[:]); err != nil {
		return PublicKey{}, err
	}
	var pubKey PublicKey
	pubKey.SetBytes(b)
	return pubKey, nil
}

func ValidPublicKey(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}

	if len(s) != 2*PubKeyLength || !isHex(s) {
		return false
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return false
	}

	if _, err := crypto.DecompressPubkey(b[:]); err != nil {
		return false
	}

	return true
}

func PubkeyToAddress(pub PublicKey) ethcmn.Address {
	pubKey, err := crypto.DecompressPubkey(pub[:])
	if err != nil {
		return ethcmn.Address{}
	}
	return crypto.PubkeyToAddress(*pubKey)
}
