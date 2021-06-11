package model

import (
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type DPosPoolLog struct {
	Id int64 `xorm:"INTEGER NOT NULL
						PRIMARY KEY AUTOINCREMENT" json:"id"`
	Height   int64     `xorm:"INTEGER NOT NULL INDEX(idx_height) COMMENT('区块高度')" json:"height"`
	Balance  string    `xorm:"TEXT NOT NULL COMMENT('矿池余额')" json:"balance"`
	Mined    string    `xorm:"TEXT NOT NULL COMMENT('挖矿量')" json:"mined"`
	Released string    `xorm:"TEXT NOT NULL COMMENT('实际释放量，实际释放量约等于挖矿量')" json:"released"`
	Total    string    `xorm:"TEXT NOT NULL COMMENT('总股权')" json:"total"`
	BlockAt  time.Time `xorm:"'blockat' DATETIME NOT NULL INDEX(idx_blockat) COMMENT('区块时间')" json:"blockAt"`
}

func (t DPosPoolLog) TableName() string {
	return "dpos_poollog"
}

type DPosTcnLog struct {
	Id int64 `xorm:"INTEGER NOT NULL
				PRIMARY KEY AUTOINCREMENT" json:"id"`
	Height    int64     `xorm:"INTEGER NOT NULL INDEX(idx_height) COMMENT('区块高度')" json:"height"`
	Address   string    `xorm:"TEXT NOT NULL INDEX(idx_address) COMMENT('节点地址')" json:"address"`
	Mortgaged string    `xorm:"TEXT NOT NULL COMMENT('节点抵押量')" json:"mortgaged"`
	Voted     string    `xorm:"TEXT NOT NULL COMMENT('用户抵押量')" json:"voted"`
	Voters    string    `xorm:"TEXT NOT NULL COMMENT('抵押用户数')" json:"voters"`
	Profit    string    `xorm:"TEXT NOT NULL COMMENT('收益')" json:"profit"`
	BlockAt   time.Time `xorm:"'blockat' DATETIME NOT NULL INDEX(idx_blockat) COMMENT('区块时间')" json:"blockAt"`
}

func (t DPosTcnLog) TableName() string {
	return "dpos_tcnlog"
}

type DPosTinLog struct {
	Id int64 `xorm:"INTEGER NOT NULL
				PRIMARY KEY AUTOINCREMENT" json:"id"`
	Height    int64     `xorm:"INTEGER NOT NULL INDEX(idx_height) COMMENT('区块高度')" json:"height"`
	Address   string    `xorm:"TEXT NOT NULL INDEX(idx_address) COMMENT('用户地址')" json:"address"`
	Validator string    `xorm:"TEXT NOT NULL COMMENT('用户选举的节点地址')" json:"validator"`
	Mortgaged string    `xorm:"TEXT NOT NULL COMMENT('用户抵押量')" json:"mortgaged"`
	Profit    string    `xorm:"TEXT NOT NULL COMMENT('收益')" json:"profit"`
	BlockAt   time.Time `xorm:"'blockat' DATETIME NOT NULL INDEX(idx_blockat) COMMENT('区块时间')" json:"blockAt"`
}

func (t DPosTinLog) TableName() string {
	return "dpos_tinlog"
}

type DPosRankLog struct {
	Id int64 `xorm:"INTEGER NOT NULL
				PRIMARY KEY AUTOINCREMENT" json:"id"`
	Height    int64     `xorm:"INTEGER NOT NULL INDEX(idx_height) COMMENT('区块高度')" json:"height"`
	Address   string    `xorm:"TEXT NOT NULL INDEX(idx_address) COMMENT('用户地址')" json:"address"`
	Mortgaged string    `xorm:"TEXT NOT NULL COMMENT('节点抵押量')" json:"mortgaged"`
	Voted     string    `xorm:"TEXT NOT NULL COMMENT('用户抵押量')" json:"voted"`
	Voters    string    `xorm:"TEXT NOT NULL COMMENT('抵押用户数')" json:"voters"`
	Total     string    `xorm:"TEXT NOT NULL COMMENT('总股权')" json:"total"`
	Rank      uint32    `xorm:"INTEGER NOT NULL COMMENT('排名')" json:"rank"`
	BlockAt   time.Time `xorm:"'blockat' DATETIME NOT NULL INDEX(idx_blockat) COMMENT('区块时间')" json:"blockAt"`
}

func (t DPosRankLog) TableName() string {
	return "dpos_ranklog"
}
