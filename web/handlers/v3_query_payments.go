package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Summary 查询所有Payment
// @Description 查询所有Payment
// @Tags v3-query
// @Accept json
// @Produce json
// @Param symbol query string false "币种"
// @Param contract query string false "币种合约地址"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Payment "成功"
// @Router /v3/payments [get]
func (hd *Handler) QueryV3Payments(ctx *gin.Context) {
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	symbol := ctx.Query("symbol")
	contract := ctx.Query("contract")
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)
	result, err := hd.dbo3.QueryV3Payments(symbol, contract, begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据用户地址查询Payment
// @Description 根据用户地址查询Payment
// @Tags v3-query
// @Accept json
// @Produce json
// @Param address path string true "账户地址"
// @Param symbol query string false "币种"
// @Param contract query string false "币种合约地址"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Payment "成功"
// @Router /v3/accounts/{address}/payments [get]
func (hd *Handler) QueryV3AccPayments(ctx *gin.Context) {
	address := ctx.Param("address")
	symbol := ctx.Query("symbol")
	contract := ctx.Query("contract")
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)

	if address == "" {
		hd.responseWrite(ctx, false, "param address is required")
		return
	}

	result, err := hd.dbo3.QueryV3AccountPayments(address, symbol, contract, begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据txhash查询Payment
// @Description 根据txhash查询Payment
// @Tags v3-query
// @Accept json
// @Produce json
// @Param txhash path string true "交易hash"
// @Param symbol query string false "币种"
// @Param contract query string false "币种合约地址"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Payment "成功"
// @Router /v3/transactions/{txhash}/payments [get]
func (hd *Handler) QueryV3TxPayments(ctx *gin.Context) {
	txhash := ctx.Param("txhash")
	symbol := ctx.Query("symbol")
	contract := ctx.Query("contract")
	order := ctx.Query("order")
	limits := ctx.Query("limit")
	cursors := ctx.Query("cursor")
	begins := ctx.Query("begin")
	ends := ctx.Query("end")
	limit, _ := strconv.ParseUint(limits, 10, 64)
	cursor, _ := strconv.ParseUint(cursors, 10, 64)
	begin, _ := strconv.ParseUint(begins, 10, 64)
	end, _ := strconv.ParseUint(ends, 10, 64)

	if txhash == "" {
		hd.responseWrite(ctx, false, "param txhash is required")
		return
	}

	result, err := hd.dbo3.QueryV3TxPayments(txhash, symbol, contract, begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}

// @Summary 根据区块高度查询Payment
// @Description 根据区块高度查询Payment
// @Tags v3-query
// @Accept json
// @Produce json
// @Param height path string true "区块高度"
// @Param symbol query string false "币种"
// @Param contract query string false "币种合约地址"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Success 200 {array}  database.V3Payment "成功"
// @Router /v3/ledgers/{height}/payments [get]
func (hd *Handler) QueryV3LedgerPayments(ctx *gin.Context) {
	heights := ctx.Param("height")
	symbol := ctx.Query("symbol")
	contract := ctx.Query("contract")
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
		hd.responseWrite(ctx, false, "param height is required")
		return
	}

	result, err := hd.dbo3.QueryV3BlockPayments(int64(height), symbol, contract, begin, end, cursor, limit, order)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
	} else {
		hd.responseWrite(ctx, true, result)
	}
}
