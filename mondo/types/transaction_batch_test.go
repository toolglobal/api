package types

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func Test_sign(t *testing.T) {
	tx := NewTxBatch()
	tx.CreatedAt = uint64(1618901550)
	tx.GasLimit = 21000 * 2
	tx.GasPrice = big.NewInt(100)
	tx.Nonce = 1
	tx.Sender, _ = HexToPubkey("03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8")

	for i := 0; i < 2; i++ {
		var op TxOp
		op.To.SetBytes(common.HexToAddress("0xc8F516fa76868b4C16bA439F3131911828339Ed5").Bytes())
		op.Value = big.NewInt(int64(i + 1))
		tx.Ops = append(tx.Ops, op)
	}
	tx.Memo = []byte("sys")
	tx.Signature, _ = tx.Sign("7fffe4e426a6772ae8a1c0f2425a90fc6320d23e416fb6d83802889fa846faa2")
	fmt.Printf("signature:%s\nhash:%s\n", hex.EncodeToString(tx.Signature), tx.Hash().Hex())
}
