package types

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_validator_encode_decode(t *testing.T) {
	v := NewValidator()
	b, _ := base64.StdEncoding.DecodeString("MkRG9w+eCnFrZxs+9Hz/Y6V3YJ3Q7j999AuFhUjQAsQ=")
	copy(v.PubKey[:], b)
	v.Nonce = 456
	v.Power = 12345
	v.Genesis = true
	v.RunFor = true
	v.Redeem = true
	v.Voted = new(big.Int)
	v.AddVoter(common.HexToAddress("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F"), new(big.Int).SetInt64(123))
	for i := 0; i < 102; i++ {
		v.AddVoter(common.BytesToAddress([]byte{uint8(i)}), new(big.Int).SetInt64(456))
	}
	v.UpdateVoter(common.HexToAddress("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F"), new(big.Int).SetInt64(456))
	fmt.Println(len(v.Voters), cap(v.Voters))

	b = v.ToBytes()
	dst := NewValidator()
	dst.FromBytes(b)
	fmt.Println(len(dst.Voters), cap(dst.Voters))
	if !v.Equal(dst) {
		t.Fatal("not equal")
	}
}
