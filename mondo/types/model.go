package types

import (
	"time"
)

type DPOSPoolLog struct {
	Id       uint64    `db:"id" json:"id"`             // 数据库自增id
	Height   uint64    `db:"height" json:"height"`     // 去看高度
	Balance  string    `db:"balance" json:"balance"`   // 矿池余额
	Mined    string    `db:"mined" json:"mined"`       // 挖矿量
	Released string    `db:"released" json:"released"` // 实际释放量，实际释放量约等于挖矿量
	Total    string    `db:"total" json:"total"`       // 总股权
	BlockAt  time.Time `db:"blockat" json:"blockAt"`   // 区块时间
}

type DPOSTcnLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 节点抵押量
	Voted     string    `db:"voted" json:"voted"`         // 用户投票量
	Voters    uint64    `db:"voters" json:"voters"`       // 投票用户数
	Profit    string    `db:"profit" json:"profit"`       // 收益
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
}

type DPOSTinLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 用户地址
	Validator string    `db:"validator" json:"validator"` // 用户选举的节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 用户抵押量
	Profit    string    `db:"profit" json:"profit"`       // 收益
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
}

type DPOSRankLog struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Height    uint64    `db:"height" json:"height"`       // 区块高度
	Address   string    `db:"address" json:"address"`     // 节点地址
	Mortgaged string    `db:"mortgaged" json:"mortgaged"` // 节点抵押量
	Voted     string    `db:"voted" json:"voted"`         // 用户投票量
	Voters    uint64    `db:"voters" json:"voters"`       // 投票用户数
	Total     string    `db:"total" json:"total"`         // 总股权
	Rank      uint32    `db:"rank" json:"rank"`           // 排名
	BlockAt   time.Time `db:"blockat" json:"blockAt"`     // 区块时间
}
