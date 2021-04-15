package types

import (
	"fmt"
	"math/big"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
)

// AppHeader ledger header for app
type AppHeader struct {
	Height     *big.Int       // refresh by new block
	ClosedAt   time.Time      // ON abci beginBlock
	BlockHash  ethcmn.Hash    // fill on new block
	Validator  ethcmn.Address // refresh by new block
	StateRoot  ethcmn.Hash    // fill after statedb commit
	XStateRoot ethcmn.Hash    // fill after xstatedb commit
	TxCount    uint64         // fill when ready to save
	PrevHash   ethcmn.Hash    // global, just used to calculate header-hash
	MinePool   *MinePool
	Params     *Config
}

func (h *AppHeader) String() string {
	return fmt.Sprintf("H:%v T:%v Val:%v R1:%v R2:%v TX:%v PREV:%v",
		h.Height, h.ClosedAt, h.Validator.Hex(), h.StateRoot.Hex(), h.XStateRoot.Hex(), h.TxCount, h.PrevHash.Hex())
}

// Hash hash
func (h *AppHeader) Hash() []byte {
	s := make([][]byte, 4)
	bz0 := h.StateRoot.Bytes()
	s[0] = make([]byte, len(bz0))
	copy(s[0], bz0)

	bz1 := h.XStateRoot.Bytes()
	s[1] = make([]byte, len(bz1))
	copy(s[1], bz1)

	bz2 := h.MinePool.ToBytes()
	s[2] = make([]byte, len(bz2))
	copy(s[2], bz2)

	bz3 := h.Params.ToBytes()
	s[3] = make([]byte, len(bz3))
	copy(s[3], bz3)

	return merkle.HashFromByteSlices(s)
}
