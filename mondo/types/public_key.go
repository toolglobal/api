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

func (pk PublicKey) Copy() PublicKey {
	p := PublicKey{}
	p.SetBytes(pk[:])
	return p
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
func (pk PublicKey) IsAddress() bool {
	for i := 0; i < 13; i++ {
		if pk[i] != 0 {
			return false
		}
	}
	return true
}

func (pk PublicKey) Bytes() []byte { return pk[:] }

func (pk PublicKey) Big() *big.Int { return new(big.Int).SetBytes(pk[:]) }

func (pk PublicKey) Hex() string {
	return hexutil.Encode(pk[:])
}

// String implements fmt.Stringer.
func (pk PublicKey) String() string {
	return pk.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (pk PublicKey) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), pk[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(pk) it will panic.
func (pk *PublicKey) SetBytes(b []byte) {
	if len(b) > len(pk) {
		b = b[len(b)-PubKeyLength:]
	}
	copy(pk[PubKeyLength-len(b):], b)
}

// MarshalText returns the hex representation of pk.
func (pk PublicKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(pk[:]).MarshalText()
}

// UnmarshalText parses pk hash in hex syntax.
func (pk *PublicKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, pk[:])
}

// UnmarshalJSON parses pk hash in hex syntax.
func (pk *PublicKey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, pk[:])
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
