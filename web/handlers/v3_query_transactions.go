package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Summary 根据txhash查询交易记录
// @Description 根据txhash查询交易记录
// @Tags v3-query
// @Accept json
// @Produce json
// @Param txhash path string true "交易hash"
// @Success 200 {array}  database.V3Transaction "成功"
// @Router /v3/transactions/{txhash} [get]
func (hd *Handler) QueryV3SingleTx(ctx *gin.Context) {
	txhash := ctx.Param("txhash")
	if txhash == "" {
		hd.responseWrite(ctx, false, "param txhash is required")
		return
	}

	result, err := hd.dbo3.QueryV3SingleTx(txhash)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据所有交易记录
// @Description 根据所有交易记录
// @Tags v3-query
// @Accept json
// @Produce json
// @Param begin query int false "开始时间戳"
// @Param end query int false "结束时间戳"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Transaction "成功"
// @Router /v3/transactions [get]
func (hd *Handler) QueryV3Txs(ctx *gin.Context) {
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)
	result, err := hd.dbo3.QueryV3Txs(begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据账户地址查询交易记录
// @Description 根据账户地址查询交易记录
// @Tags v3-query
// @Accept json
// @Produce json
// @Param address path string true "账户地址"
// @Param begin query int false "开始时间戳"
// @Param end query int false "结束时间戳"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Transaction "成功"
// @Router /v3/accounts/{address}/transactions [get]
func (hd *Handler) QueryV3AccTxs(ctx *gin.Context) {
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)
	address := ctx.Param("address")

	if address == "" {
		hd.responseWrite(ctx, false, "param address is required")
		return
	}

	result, err := hd.dbo3.QueryV3AccountTxs(address, begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据区块高度查询交易记录
// @Description 根据区块高度查询交易记录
// @Tags v3-query
// @Accept json
// @Produce json
// @Param height path string true "区块高度"
// @Param begin query int false "开始时间戳"
// @Param end query int false "结束时间戳"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Transaction "成功"
// @Router /v3/ledgers/{height}/transactions [get]
func (hd *Handler) QueryV3LedgerTxs(ctx *gin.Context) {
	heights := ctx.Param("height")
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	height, _ := strconv.ParseUint(heights, 10, 64)
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)

	if height == 0 {
		hd.responseWrite(ctx, false, "param address is required")
		return
	}

	result, err := hd.dbo3.QueryV3BlockTxs(int64(height), begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}
