package bean

type Ticker struct {
	Amount string  `json:"amount"` // 交易金额
	High   string  `json:"high"`   // 最高价
	Vol    string  `json:"vol"`    // 成交量
	Last   float64 `json:"last"`   // 最新成交价
	Low    string  `json:"low"`    // 最低价
	Buy    float64 `json:"buy"`    // 买入价
	Sell   float64 `json:"sell"`   // 卖出价
	Rose   string  `json:"rose"`   // 振幅
	Time   int64   `json:"time"`   // 时间戳 ms
}

type EbuycoinResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Ticker `json:"data"`
}
