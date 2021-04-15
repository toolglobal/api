package bean

import "github.com/ethereum/go-ethereum/common"

type PublicResp struct {
	IsSuccess bool        `json:"isSuccess"` // 是否成功
	Message   string      `json:"message"`   // 消息提醒
	Result    interface{} `json:"result"`    // 数据对象
}

type V2GenKeyResult struct {
	Privkey string `json:"privkey"` // 私钥
	Pubkey  string `json:"pubkey"`  // 公钥（压缩公钥）
	Address string `json:"address"` // 地址
}

type V2AccountResult struct {
	Address string `json:"address"` // 地址
	Balance string `json:"balance"` // 余额
	Nonce   uint64 `json:"nonce"`   // nonce
}

type V2ConvertResult struct {
	OldAddress string `json:"old_address"` // 旧地址，原公钥
	PublicKey  string `json:"public_key"`  // 公钥，同上
	Address    string `json:"address"`     // 地址
}

type V2ContractActResult struct {
	Address  string `json:"address"`  // 地址
	Balance  string `json:"balance"`  // 余额
	Nonce    uint64 `json:"nonce"`    // nonce
	Code     string `json:"code"`     // 合约字节码
	Suicided bool   `json:"suicided"` // 合约是否已自杀
}

// EVM事件日志
type V2EvmLog struct {
	Address     common.Address `json:"address"`          // 合约地址
	Topics      []common.Hash  `json:"topics"`           // 事件主题
	Data        []byte         `json:"data"`             // 数据
	BlockNumber uint64         `json:"blockNumber"`      // 区块高度
	TxHash      common.Hash    `json:"transactionHash"`  // 交易hash
	TxIndex     uint           `json:"transactionIndex"` // 交易索引
	BlockHash   common.Hash    `json:"blockHash"`        // 区块hash
	Index       uint           `json:"logIndex"`         // 日志索引
	Removed     bool           `json:"removed"`          // 是否已移除
}
