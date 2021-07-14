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
	Height    *big.Int       // refresh by new block
	ClosedAt  time.Time      // ON abci beginBlock
	BlockHash ethcmn.Hash    // fill on new block
	Validator ethcmn.Address // refresh by new block
	MinerFee  *big.Int       // miner fee, added on deliverTx
	StateRoot ethcmn.Hash    // fill after statedb commit
	TxCount   uint64         // fill when ready to save
	PrevHash  ethcmn.Hash    // global, just used to calculate header-hash
}

func (h *AppHeader) Copy() *AppHeader {
	header := &AppHeader{
		Height:    new(big.Int),
		ClosedAt:  h.ClosedAt,
		BlockHash: h.BlockHash,
		Validator: h.Validator,
		MinerFee:  h.MinerFee,
		StateRoot: h.StateRoot,
		TxCount:   h.TxCount,
		PrevHash:  h.PrevHash,
	}

	return header
}

func (h *AppHeader) String() string {
	return fmt.Sprintf("H:%v T:%v Val:%v R:%v TX:%v PREV:%v Fee:%v",
		h.Height, h.ClosedAt, h.Validator.Hex(), h.StateRoot.Hex(), h.TxCount, h.PrevHash.Hex(), h.MinerFee.String())
}

// Hash hash
func (h *AppHeader) Hash() []byte {
	s := make([][]byte, 1)
	bz0 := h.StateRoot.Bytes()
	s[0] = make([]byte, len(bz0))
	copy(s[0], bz0)

	return merkle.HashFromByteSlices(s)
}
