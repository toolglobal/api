package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/big"
	"reflect"
)

type Validator struct {
	PubKey     ed25519.PubKey // 节点公钥
	Nonce      uint64         // nonce
	Power      uint64         // 投票权
	Genesis    bool           // 是否是创世节点
	RunFor     bool           // 参与竞选
	Balance    *big.Int       // 余额（累积收益）
	Profit     *big.Int       // 本轮收益
	Mortgaged  *big.Int       // 已抵押
	ToMortgage *big.Int       // 待抵押
	Redeem     bool           // 待赎回
	Voters     []Vot          // 选民列表
	Voted      *big.Int       // 选民贡献
}

func NewValidator() *Validator {
	return &Validator{
		PubKey:     make([]byte, ed25519.PubKeySize),
		Balance:    big.NewInt(0),
		Profit:     big.NewInt(0),
		Mortgaged:  big.NewInt(0),
		ToMortgage: big.NewInt(0),
		Voted:      big.NewInt(0),
		Voters:     make([]Vot, 0, 100),
	}
}

type Vot struct {
	Address common.Address // 选民地址
	Amount  *big.Int       // 投票金额
}

func (v *Validator) String() string {
	b, _ := v.ToJSON()
	return string(b)
}

func (v *Validator) Address() crypto.Address {
	return v.PubKey.Address()
}

func (v *Validator) PublicKeyStr() string {
	return v.PubKey.String()
}

func (v *Validator) ToBytes() []byte {
	b, err := rlp.EncodeToBytes(v)
	if err != nil {
		panic(err)
	}
	return b
}

func (v *Validator) FromBytes(b []byte) {
	err := rlp.DecodeBytes(b, v)
	if err != nil {
		panic(err)
	}
}

func (v *Validator) ToJSON() ([]byte, error) {
	return json.Marshal(&v)
}

func (v *Validator) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(&v, "", "    ")
}

func (v *Validator) FromJSON(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *Validator) Equal(to *Validator) bool {
	return reflect.DeepEqual(v, to)
}

func (v *Validator) TotalShare() *big.Int {
	return new(big.Int).Add(v.Mortgaged, v.Voted)
}

func (v *Validator) TCNShare() *big.Int {
	return v.Mortgaged
}

func (v *Validator) TINShare() *big.Int {
	return v.Voted
}

func (v *Validator) UpdateVoter(address common.Address, amount *big.Int) {
	for i, voter := range v.Voters {
		if voter.Address == address {
			v.Voted = new(big.Int).Sub(v.Voted, voter.Amount)
			v.Voted = new(big.Int).Add(v.Voted, amount)
			v.Voters[i].Amount = new(big.Int).Set(amount)
			return
		}
	}
}

func (v *Validator) AddVoter(address common.Address, amount *big.Int) {
	v.Voted = new(big.Int).Add(v.Voted, amount)

	voter := Vot{
		Address: address,
		Amount:  new(big.Int).Set(amount),
	}
	len, cap := len(v.Voters), cap(v.Voters)
	if len == cap {
		ns := make([]Vot, len+1, len+100) //增加增量
		copy(ns[:], v.Voters)
		ns[len] = voter
		v.Voters = ns
	} else {
		v.Voters = append(v.Voters, voter)
	}
}

func (v *Validator) RemoveVoter(address common.Address) {
	for idx, voter := range v.Voters {
		if voter.Address == address {
			v.Voted = new(big.Int).Sub(v.Voted, voter.Amount)
			ns := v.Voters[:len(v.Voters)-1]
			copy(ns[idx:], v.Voters[idx+1:])
			v.Voters = ns
			break
		}
	}
}
