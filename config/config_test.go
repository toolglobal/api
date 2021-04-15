package config

import (
	"fmt"
	"testing"
)

func Test_parsefile(t *testing.T) {
	cfg := New()
	if err := cfg.Init("./config.toml"); err != nil {
		panic("On init yaml:" + err.Error())
	}
	fmt.Println(cfg)
}
