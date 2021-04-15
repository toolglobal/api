package handlers

import (
	"github.com/axengine/httpc"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/web/bean"
	"strings"
	"time"
)

// @Summary 查询币种价格
// @Description 查询币种价格：币种-USDT交易对
// @Tags v3-extension
// @Accept json
// @Produce json
// @Param symbol path string false "币种"
// @Success 200 {object}  bean.Ticker "成功"
// @Router /v3/ext/price/{symbol} [get]
func (hd *Handler) V3QueryPrice(ctx *gin.Context) {
	symbol := ctx.Param("symbol") + "usdt"
	var resp bean.EbuycoinResponse

	ctxx := httpc.NewContext()
	ctxx.SetTotalTimeout(time.Second * 10)

	err := httpc.New("https://openapi.ebuycoin.com").SetContext(ctxx).Path("open").Path("api").
		Path("get_ticker").Query("symbol", strings.ToLower(symbol)).Get(&resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.responseWrite(ctx, true, resp.Data)
}
