package types

type BeanEvmCallResult = struct {
	Code    uint32 `json:"code"`    // 错误码
	Msg     string `json:"msg"`     // msg
	Ret     string `json:"ret"`     // 返回数据的hex编码
	GasUsed uint64 `json:"gasUsed"` // 消耗的gas
}

type BeanContract struct {
	Address  string `json:"address"`  // 地址
	Balance  string `json:"balance"`  // 余额
	Nonce    uint64 `json:"nonce"`    // nonce
	Code     string `json:"code"`     // 合约代码
	Suicided bool   `json:"suicided"` // 合约是否已自杀
}

type BeanAccount struct {
	Address string `json:"address"` // 地址
	Balance string `json:"balance"` // 余额
	Nonce   uint64 `json:"nonce"`   // nonce
}
