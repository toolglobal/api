package types

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"
)

var (
	USER_ADDRESS = "0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F"
	USER_PUBKEY  = "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8"
	USER_PRIVKEY = "7fffe4e426a6772ae8a1c0f2425a90fc6320d23e416fb6d83802889fa846faa2"
	NODE_ADDRESS = "5B817CD2101208A68461AEF4D9635AA512FD70B7"
)

func TestTxUserDelegate_sign(t *testing.T) {
	tx := NewUserDelegateTx()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, _ := hex.DecodeString(USER_PUBKEY)
	copy(tx.Sender[:], b)
	tx.Nonce = 1
	tx.OpType = USER_OPTYPE_MORTGAGE
	tx.OpValue = big.NewInt(10000 * 1e8)
	b, _ = hex.DecodeString(NODE_ADDRESS)
	tx.Receiver = b

	sign, err := tx.Sign(USER_PRIVKEY)
	if err != nil {
		t.Fatal(err)
	}
	tx.Signature = sign

	fmt.Println(sign)

	if !tx.Verify() {
		t.Fatal("verify failed")
	}
}
