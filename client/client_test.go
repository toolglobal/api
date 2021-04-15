package client

import (
	"context"
	"github.com/wolot/api/config"
	"testing"

	"github.com/wolot/api/database"
	"github.com/wolot/api/database/basesql"
	"github.com/wolot/api/datamanager"
	"github.com/wolot/api/libs/log"
)

var dataMgr *datamanager.DataManager

func init() {
	dataM, err := datamanager.NewDataManager("mondo_query.db", func(dbname string) database.Database {
		dbi := &basesql.Basesql{}
		err := dbi.Init(dbname, "", log.Logger)
		if err != nil {
			panic(err)
		}
		return dbi
	})

	if err != nil {
		panic(err)
	}

	dataMgr = dataM
}

func TestFetch_FetchBlockInfo(t *testing.T) {
	height := int64(1)
	fetch := NewFetch("http://192.168.8.101:26657")
	block, err := fetch.FetchBlockInfo(height)
	if err != nil {
		t.Fatal(err)
	}

	if block.Block.Header.Height != height {
		t.Fatal(block.Block.Header.Height)
	}
}

func TestClient(t *testing.T) {
	client, err := NewClient(context.Background(), config.Tokens{}, 3, "http://192.168.8.123:26657", dataMgr, 0)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer dataMgr.Close()

	client.Start()
}
