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
	ChainId     string
	Versions    []int
	StartHeight int64
	TGSBaseURL  string
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

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
