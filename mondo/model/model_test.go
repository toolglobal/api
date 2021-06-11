package model

import "testing"
import (
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

func Test_sync(t *testing.T) {
	engine, _ := xorm.NewEngine("sqlite3", "test.db")
	engine.Sync2(new(DPosPoolLog),
		new(DPosTcnLog),
		new(DPosTinLog),
		new(DPosRankLog))
}
