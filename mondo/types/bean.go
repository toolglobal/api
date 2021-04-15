package types

import "encoding/json"

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

type BeanValidator struct {
	Address    string `json:"address"`    // 节点公钥
	Nonce      uint64 `json:"nonce"`      // nonce
	Power      uint64 `json:"power"`      // 投票权
	Genesis    bool   `json:"genesis"`    // 是否是创世节点
	RunFor     bool   `json:"runFor"`     // 参与竞选
	Balance    string `json:"balance"`    // 余额（累积收益）
	Profit     string `json:"-"`          // 本轮收益
	Mortgaged  string `json:"mortgaged"`  // 已抵押
	ToMortgage string `json:"toMortgage"` // 待抵押
	Redeem     bool   `json:"redeem"`     // 待赎回
	Voters     int    `json:"voters"`     // 投票人数
	Voted      string `json:"voted"`      // 选民贡献
}

func (v *BeanValidator) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(&v, "", "    ")
}

type BeanVots struct {
	Address string    `json:"address"` // 节点公钥
	Voted   string    `json:"voted"`   // 选民贡献
	Voters  []BeanVot `json:"voters"`  // 选民列表
}

func (v *BeanVots) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(&v, "", "    ")
}

type BeanVot struct {
	Address string `json:"address"` // 选民地址
	Amount  string `json:"amount"`  // 投票金额
}

type BeanVoter struct {
	Address    string `json:"address"`    // 选民地址
	Balance    string `json:"balance"`    // 累积收益余额
	Validator  string `json:"validator"`  // 节点地址
	Mortgaged  string `json:"mortgaged"`  // 已抵押
	ToMortgage string `json:"toMortgage"` // 待抵押
	Burned     string `json:"burned"`     // 已燃烧
	Redeem     bool   `json:"redeem"`     // 赎回标志
}

type BeanDPOSPool struct {
	CurrentHeight int64 `json:"currentHeight"` // 当前高度
	LastHeight    int64 `json:"lastHeight"`    // 上一轮DPOS高度
	NextHeight    int64 `json:"nextHeight"`    // 下一轮DPOS高度
}
