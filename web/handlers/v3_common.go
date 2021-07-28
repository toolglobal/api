package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/toolglobal/api/client"
)

// @Summary 查询已发行代币
// @Description 查询已发行代币
// @Tags v3-config
// @Accept json
// @Produce json
// @Success 200 {object}  config.Tokens "成功"
// @Router /v3/config/tokens [get]
func (hd *Handler) V3QueryConfigTokens(ctx *gin.Context) {
	cli := client.NewTokenMgr(hd.cfg.TGSBaseURL, hd.cfg.ChainId)
	if err := cli.Sync(); err != nil {
		hd.responseWrite(ctx, false, nil)
		return
	}
	hd.responseWrite(ctx, true, cli.Tokens())
}

// @Summary 查询节点配置信息
// @Description 查询节点配置信息
// @Tags v3-config
// @Accept json
// @Produce json
// @Success 200 {object}  config.Nodes "成功"
// @Router /v3/config/nodes [get]
//func (hd *Handler) V3QueryConfigNodes(ctx *gin.Context) {
//	hd.responseWrite(ctx, true, hd.cfg.Nodes)
//}
