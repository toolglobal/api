package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func Test_rlp_encode_decode(t *testing.T) {
	var voter Voter
	voter.Address = common.HexToAddress("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F")
	voter.Balance = new(big.Int)
	voter.Validator = common.HexToAddress("5B817CD2101208A68461AEF4D9635AA512FD70B7")
	voter.Mortgaged = new(big.Int).SetInt64(456)
	voter.ToMortgage = new(big.Int).SetInt64(789)
	voter.Burned = new(big.Int).SetInt64(1024)
	voter.Redeem = true

	b := voter.ToBytes()
	var dst Voter
	dst.FromBytes(b)
	fmt.Println(dst)

	if !voter.Equal(&dst) {
		t.Fatal("not equal")
	}
}
