package database

import (
	"math/big"
	"time"
)

type V3Ledger struct {
	Id         uint64    `db:"id" json:"id"`               // 数据库自增id
	Height     int64     `db:"height" json:"height"`       // 区块高度
	BlockHash  string    `db:"blockHash" json:"blockHash"` // 区块hash
	BlockSize  int       `db:"blockSize" json:"blockSize"` // 区块大小：字节
	Validator  string    `db:"validator" json:"validator"` // 区块验证者节点地址
	TxCount    int64     `db:"txCount" json:"txCount"`     // 区块交易数
	GasLimit   int64     `db:"gasLimit" json:"gasLimit"`   // 区块gas限额之和
	GasUsed    int64     `db:"gasUsed" json:"gasUsed"`     // 区块所有交易消耗gas之和
	GasPrice   string    `db:"gasPrice" json:"gasPrice"`   // 区块交易平均gas价格，可能是小数
	CreatedAt  time.Time `db:"createdAt" json:"createdAt"` // 区块时间
	TotalPrice *big.Int  `db:"-" json:"-"`                 //
}

type V3Transaction struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Hash      string    `db:"hash" json:"hash"`           // 交易hash
	Height    int64     `db:"height" json:"height"`       // 区块高度
	Typei     int       `db:"typei" json:"typei"`         // 交易类型
	Types     string    `db:"types" json:"types"`         // 交易类型
	Sender    string    `db:"sender" json:"sender"`       // 交易发起者地址
	Nonce     int64     `db:"nonce" json:"nonce"`         // 交易发起者nonce
	Receiver  string    `db:"receiver" json:"receiver"`   // 交易接受者地址
	Value     string    `db:"value" json:"value"`         // 交易金额
	GasLimit  int64     `db:"gasLimit" json:"gasLimit"`   // gas限额
	GasUsed   int64     `db:"gasUsed" json:"gasUsed"`     // gas使用量
	GasPrice  string    `db:"gasPrice" json:"gasPrice"`   // gas价格
	Memo      string    `db:"memo" json:"memo"`           // 备注
	Payload   string    `db:"payload" json:"payload"`     // 负载
	Events    string    `db:"events" json:"events"`       // 交易事件
	Codei     uint32    `db:"codei" json:"codei"`         // 失败代码
	Codes     string    `db:"codes" json:"codes"`         // 失败原因
	CreatedAt time.Time `db:"createdAt" json:"createdAt"` // 区块时间
}

type V3Payment struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Hash      string    `db:"hash" json:"hash"`           // 交易hash
	Height    int64     `db:"height" json:"height"`       // 区块高度
	EvName    string    `db:"evName" json:"evName"`       // 事件名称
	Idx       uint      `db:"idx" json:"idx"`             // 交易索引
	Sender    string    `db:"sender" json:"sender"`       // 转账发起方地址
	Receiver  string    `db:"receiver" json:"receiver"`   // 转账接受方地址
	Symbol    string    `db:"symbol" json:"symbol"`       // 币种，原生币为“OLO”
	Contract  string    `db:"contract" json:"contract"`   // 合约地址，原生币为空或全零黑洞地址
	Value     string    `db:"value" json:"value"`         // 交易金额
	CreatedAt time.Time `db:"createdAt" json:"createdAt"` // 区块时间
}
