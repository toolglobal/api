package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/toolglobal/api/config"
	"github.com/toolglobal/api/web/dbo"
	"go.uber.org/zap"
	"sync"
)

type Handler struct {
	client *http.HTTP
	logger *zap.Logger
	mu     sync.Mutex
	cfg    *config.Config
	dbo3   *dbo.DBO
}

func NewHandler(logger *zap.Logger, cfg *config.Config, dbo3 *dbo.DBO) *Handler {
	var h Handler
	h.client, _ = http.New("http://"+cfg.RPC, "/websocket")
	h.logger = logger
	h.cfg = cfg
	h.dbo3 = dbo3
	return &h
}

func (hd *Handler) responseWrite(ctx *gin.Context, isSuccess bool, result interface{}) {
	ret := gin.H{
		"isSuccess": isSuccess,
	}

	if isSuccess {
		ret["result"] = result
	} else {
		ret["message"] = result
	}

	ctx.JSON(200, ret)
}

func (hd *Handler) responseWriteV2(ctx *gin.Context, isSuccess bool, result interface{}, msg string) {
	ret := gin.H{
		"isSuccess": isSuccess,
	}
	ret["result"] = result
	if !isSuccess {
		ret["message"] = msg
	} else {
		ret["message"] = "success"
	}
	ctx.JSON(200, ret)
}
