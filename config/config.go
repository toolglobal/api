package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	Bind        string
	RPC         string
	Dev         bool
	Metrics     bool
	Versions    []int
	StartHeight int64
	Tokens      Tokens
	Nodes       Nodes
	Limiter     Limiter
}

func New() *Config {
	return &Config{}
}

func (p *Config) Init(cfgFile string) error {
	_, err := toml.DecodeFile(cfgFile, p)
	return err
}

type Limiter struct {
	Interval duration
	Capacity int64
}

type LogInfo struct {
	Path string
}

type Tokens struct {
	Coins []ERC20
}

type ERC20 struct {
	Name     string `json:"name"`     // 代币名称，例如Tool Global Blockchain
	Symbol   string `json:"symbol"`   // 代币symbol，例如OLO
	Decimals int    `json:"decimals"` // 小数位数
	//TotalSupply string `json:"totalSupply"` // 总发行量，整数
	//CreatedAt   string `json:"createdAt"`   // 发行时间
	Address string `json:"address"` // 合约地址
	//Exchange    bool   `json:"exchange"`    // 是否在交易所交易
	//WhitePaper  string `json:"whitePaper"`  // 白皮书
	//Website     string `json:"website"`     // 官网
	Icon string `json:"icon"` // icon url
}

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Address string `json:"address"` // 节点地址
	Name    string `json:"name"`    // 节点名称
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
