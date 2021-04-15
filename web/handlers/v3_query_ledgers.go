package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Summary 查询所有账本
// @Description 查询所有账本
// @Tags v3-query
// @Accept json
// @Produce json
// @Param begin query int false "开始时间戳"
// @Param end query int false "结束时间戳"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Ledger "成功"
// @Router /v3/ledgers [get]
func (hd *Handler) QueryV3Ledgers(ctx *gin.Context) {
	order := ctx.Query("order")
	limit := ctx.Query("limit")
	cursor := ctx.Query("cursor")
	begin := ctx.Query("begin")
	end := ctx.Query("end")

	iLimit, _ := strconv.ParseUint(limit, 10, 64)
	iCursor, _ := strconv.ParseUint(cursor, 10, 64)
	iBegin, _ := strconv.ParseUint(begin, 10, 64)
	iEnd, _ := strconv.ParseUint(end, 10, 64)

	result, err := hd.dbo3.QueryV3Ledgers(iBegin, iEnd, iCursor, iLimit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据高度查询账本
// @Description 根据高度查询账本
// @Tags v3-query
// @Accept json
// @Produce json
// @Param height path string true "区块高度"
// @Success 200 {array}  database.V3Ledger "成功"
// @Router /v3/ledgers/{height} [get]
func (hd *Handler) QueryV3Ledger(ctx *gin.Context) {
	height, _ := strconv.ParseUint(ctx.Param("height"), 10, 0)
	result, err := hd.dbo3.QueryV3LedgerByHeight(int64(height))
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}
