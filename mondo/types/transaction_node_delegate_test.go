package types

import (
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
	"time"
)

var (
	//NODE_ADDRESS  = "5B817CD2101208A68461AEF4D9635AA512FD70B7"
	NODE_PUB_KEY  = "2BWneDxGm6RnGdGezKbncFuOrRDljUGgDDffj714FaE="
	NODE_PRIV_KEY = "haCuQyYYgHI5WFRU0pAxo1beJ542NO0hEqsRdatGPqXYFad4PEabpGcZ0Z7MpudwW46tEOWNQaAMN9+PvXgVoQ=="
)

func TestTxNodeDelegate_sign(t *testing.T) {
	tx := NewNodeDelegateTx()
	tx.CreatedAt = uint64(time.Now().Unix())
	b, _ := base64.StdEncoding.DecodeString(NODE_PUB_KEY)
	copy(tx.Sender[:], b)
	tx.Nonce = 1
	tx.OpType = NODE_OPTYPE_MORTGAGE
	tx.OpValue = big.NewInt(10000 * 1e8)
	tx.Receiver = common.HexToAddress("0x0F508F143E77b39F8e20DD9d2C1e515f0f527D9F").Bytes()
	b, _ = base64.StdEncoding.DecodeString(NODE_PRIV_KEY)
	tx.Signature = tx.Sign(b)

	if !tx.Verify() {
		t.Fatal("verify failed")
	}
}
