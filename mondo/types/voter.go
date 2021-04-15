package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"reflect"
)

type Voter struct {
	Address    common.Address // 选民地址
	Balance    *big.Int       // 累积收益余额
	Validator  common.Address // 节点地址 使用common.Address替换crypto.Address
	Mortgaged  *big.Int       // 已抵押
	ToMortgage *big.Int       // 待抵押
	Burned     *big.Int       // 已燃烧
	Redeem     bool           // 赎回标志
}

func NewVoter() *Voter {
	return &Voter{
		Balance:    new(big.Int),
		Burned:     new(big.Int),
		Mortgaged:  new(big.Int),
		ToMortgage: new(big.Int),
	}
}

func (v *Voter) CanRemove() bool {
	var big0 = big.NewInt(0)
	return v.Balance.Cmp(big0) == 0 && v.Mortgaged.Cmp(big0) == 0 &&
		v.ToMortgage.Cmp(big0) == 0 && v.Burned.Cmp(big0) == 0 &&
		v.Redeem == false
}

func (v *Voter) String() string {
	b, _ := v.ToJSON()
	return string(b)
}

func (v *Voter) ToBytes() []byte {
	b, err := rlp.EncodeToBytes(v)
	if err != nil {
		panic(err)
	}
	return b
}

func (v *Voter) FromBytes(b []byte) {
	err := rlp.DecodeBytes(b, v)
	if err != nil {
		panic(err)
	}
}

func (v *Voter) ToJSON() ([]byte, error) {
	return json.Marshal(&v)
}

func (v *Voter) FromJSON(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *Voter) Equal(to *Voter) bool {
	return reflect.DeepEqual(v, to)
}

// 实际股权贡献
func (v *Voter) TotalShare() *big.Int {
	//// 燃烧金额*15%
	//burnShare := new(big.Int).Mul(v.Burned, new(big.Int).SetInt64(15))
	//burnShare = new(big.Int).Div(burnShare, new(big.Int).SetInt64(100))
	//
	//// 实际抵押金额
	//mortgaged := new(big.Int).Sub(v.Mortgaged, v.Burned)
	//return new(big.Int).Add(mortgaged, burnShare)
	return new(big.Int).Set(v.Mortgaged)
}

// 实际抵押金额
func (v *Voter) RealMortgaged() *big.Int {
	return new(big.Int).Sub(v.Mortgaged, v.Burned)
}
