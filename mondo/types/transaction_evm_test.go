package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func TestTxEvm_SignAndVerify(t *testing.T) {
	tx := NewTxEvm()
	tx.GasPrice = big.NewInt(1)
	tx.GasLimit = 21000
	tx.Nonce = 173
	tx.Sender, _ = HexToPubkey("0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
	to := common.HexToAddress("0xB944aC8c6E20475CA528854C29350Df7daF9d1A5")
	copy(tx.Body.To[13:], to.Bytes())
	tx.Body.Value = big.NewInt(499999997900000)
	tx.Body.Load = common.Hex2Bytes("")
	tx.Body.Memo = []byte("e5efc346-5b74-40b3-a82f-f70b8a0ae1f9")
	tx.Signature, _ = tx.Sign("a8971729fbc199fb3459529cebcd8704791fc699d88ac89284f23ff8e7fca7d6")
	fmt.Println(" len of signature=", len(tx.Signature))
	fmt.Println("signature:", hex.EncodeToString(tx.Signature))
	b := tx.Verify()
	if !b {
		t.Fatal("verify ", b)
	}
}

func TestTxEvm_rlp(t *testing.T) {
	tx := NewTxEvm()
	tx.GasPrice = big.NewInt(1)
	tx.GasLimit = 21000
	tx.Nonce = 1
	tx.Sender, _ = HexToPubkey("0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
	tx.Body.To = PublicKey{}
	tx.Body.Value = big.NewInt(0)
	tx.Body.Load = []byte("hello world")
	tx.Body.Memo = []byte("hello")
	tx.Signature, _ = tx.Sign("e4b96ad1a38d52c959ea2b4a17c1ac70771367d270f7dd26b636749fcf6d9576")
	bs := tx.ToBytes()
	tx2 := NewTxEvm()
	tx2.FromBytes(bs)
	if !bytes.Equal(bs, tx2.ToBytes()) {
		t.Fatal("not equal")
	}
}
