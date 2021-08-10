package main

import (
	"context"
	"os"

	"github.com/toolglobal/api/client"
	"github.com/toolglobal/api/config"
	"github.com/toolglobal/api/database"
	"github.com/toolglobal/api/database/basesql"
	"github.com/toolglobal/api/datamanager"
	"github.com/toolglobal/api/libs/log"
	"github.com/toolglobal/api/web/dbo"
	"github.com/toolglobal/api/web/server"
	"go.uber.org/zap"
)

var (
	BUILD_TIME string
	GIT_HASH   string
	GO_VERSION string
)

func main() {
	log.Logger.Info("init", zap.String("build", BUILD_TIME), zap.String("commit", GIT_HASH), zap.String("go", GO_VERSION))
	cfg := config.New()
	if err := cfg.Init("./config/config.toml"); err != nil {
		panic("On init toml:" + err.Error())
	}
	log.Logger.Info("config", zap.Any("cfg", cfg))

	dataM3, err := datamanager.NewDataManager("mondo_query_v3.db", func(dbname string) database.Database {
		dbi := &basesql.Basesql{}
		_ = os.Mkdir("data", 755)
		err := dbi.Init(dbname, "data", log.Logger)
		if err != nil {
			panic(err)
		}
		return dbi
	})
	if err != nil {
		panic(err)
	}
	defer dataM3.Close()

	for _, version := range cfg.Versions {
		if version == 3 {
			syncCli, err := client.NewClient(context.Background(), cfg.TGSBaseURL, cfg.ChainId, version, "http://"+cfg.RPC, dataM3, cfg.StartHeight)
			if err != nil {
				panic(err)
			}
			go syncCli.Start()
		}
	}

	server := server.NewServer(log.Logger, cfg, dbo.New(dataM3))
	server.Start()
}
