package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	API_V2_QUERY_ACCOUNT          = "/v2/account"
	API_V2_CONTRACT_QUERY_ACCOUNT = "/v2/contract/account"
	API_V2_CONTRACT_QUERY_LOGS    = "/v2/contract/logs"
	API_V2_CONTRACT_CALL          = "/v2/contract/call"
)

type QueryBase struct {
	Order  string
	Limit  uint64
	Cursor uint64

	Begin uint64
	End   uint64
}

type DPosQuery struct {
	QueryBase
	Height  uint64
	Address ethcmn.Address
}

func (q *DPosQuery) ToBytes() []byte {
	bz, err := rlp.EncodeToBytes(q)
	if err != nil {
		panic(err)
	}
	return bz
}

func (q *DPosQuery) FromBytes(bz []byte) error {
	return rlp.DecodeBytes(bz, q)
}

func (q *DPosQuery) GetQuery(bz []byte) error {
	if err := q.FromBytes(bz); err != nil {
		return err
	}
	if q.Order == "" {
		q.Order = "DESC"
	}
	if q.Limit == 0 {
		q.Limit = 10
	}

	return nil
}
