package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
	"strconv"
)

func getDPOSQuery(ctx *gin.Context) (*types.DPOSQuery, error) {
	order := ctx.Query("order")
	limit := ctx.Query("limit")
	cursor := ctx.Query("cursor")
	begin := ctx.Query("begin")
	end := ctx.Query("end")
	address := ctx.Query("address")
	height := ctx.Query("height")
	if len(address) > 0 && !types.ValidAddress(address) {
		return nil, errors.New("invalid address")
	}
	var query types.DPOSQuery
	query.Order = order
	query.Limit, _ = strconv.ParseUint(limit, 10, 64)
	query.Cursor, _ = strconv.ParseUint(cursor, 10, 64)
	query.Begin, _ = strconv.ParseUint(begin, 10, 64)
	query.End, _ = strconv.ParseUint(end, 10, 64)
	query.Address = common.HexToAddress(address)
	query.Height, _ = strconv.ParseUint(height, 10, 64)
	return &query, nil
}

// @Summary 查询矿池信息
// @Description 查询矿池信息
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Success 200 {object}  types.BeanDPOSPool "成功"
// @Router /v3/dpos/pool [get]
func (hd *Handler) V3QueryPool(ctx *gin.Context) {
	result, err := hd.client.ABCIQuery(ctx, "/v3/dpos/pool", nil)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var resp types.Result
	err = rlp.DecodeBytes(result.Response.Value, &resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if resp.Code != types.CodeType_OK {
		hd.responseWrite(ctx, false, resp.Log)
		return
	}

	var data = types.BeanDPOSPool{}
	_ = json.Unmarshal(resp.Data, &data)

	hd.responseWrite(ctx, true, data)
}

// @Summary 查询矿池记录
// @Description 查询矿池记录
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param height query int false "高度"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Success 200 {array}  types.DPOSPoolLog "成功"
// @Router /v3/dpos/poollogs [get]
func (hd *Handler) V3QueryPoolLogs(ctx *gin.Context) {
	query, err := getDPOSQuery(ctx)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_DPOS_POOLLOG, query.ToBytes())
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var resp types.Result
	err = rlp.DecodeBytes(result.Response.Value, &resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if resp.Code != types.CodeType_OK {
		hd.responseWrite(ctx, false, resp.Log)
		return
	}

	var data []types.DPOSPoolLog
	_ = json.Unmarshal(resp.Data, &data)

	hd.responseWrite(ctx, true, data)
}

// @Summary 查询节点dpos记录
// @Description 查询节点dpos记录
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address query string false "地址"
// @Param height query int false "高度"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Success 200 {array}  types.DPOSTcnLog "成功"
// @Router /v3/dpos/tcnlogs [get]
func (hd *Handler) V3QueryTcnLogs(ctx *gin.Context) {
	query, err := getDPOSQuery(ctx)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_DPOS_TCNLOG, query.ToBytes())
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var resp types.Result
	err = rlp.DecodeBytes(result.Response.Value, &resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if resp.Code != types.CodeType_OK {
		hd.responseWrite(ctx, false, resp.Log)
		return
	}

	var data []types.DPOSTcnLog
	_ = json.Unmarshal(resp.Data, &data)

	hd.responseWrite(ctx, true, data)
}

// @Summary 查询用户dpos记录
// @Description 查询用户dpos记录
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address query string false "地址"
// @Param height query int false "高度"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Success 200 {array}  types.DPOSTinLog "成功"
// @Router /v3/dpos/tinlogs [get]
func (hd *Handler) V3QueryTinLogs(ctx *gin.Context) {
	query, err := getDPOSQuery(ctx)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_DPOS_TINLOG, query.ToBytes())
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var resp types.Result
	err = rlp.DecodeBytes(result.Response.Value, &resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if resp.Code != types.CodeType_OK {
		hd.responseWrite(ctx, false, resp.Log)
		return
	}

	var data []types.DPOSTinLog
	_ = json.Unmarshal(resp.Data, &data)

	hd.responseWrite(ctx, true, data)
}

// @Summary 查询DPOS排名
// @Description 查询DPOS排名
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address query string false "地址"
// @Param height query int false "高度"
// @Param cursor query int false "游标"
// @Param limit query int false "限制"
// @Param order query string false "排序(ASC/DESC)"
// @Param begin query int false "开始时间"
// @Param end query int false "结束时间"
// @Success 200 {array}  types.DPOSRankLog "成功"
// @Router /v3/dpos/ranklogs [get]
func (hd *Handler) V3QueryRankLogs(ctx *gin.Context) {
	query, err := getDPOSQuery(ctx)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_DPOS_RANKLOG, query.ToBytes())
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var resp types.Result
	err = rlp.DecodeBytes(result.Response.Value, &resp)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if resp.Code != types.CodeType_OK {
		hd.responseWrite(ctx, false, resp.Log)
		return
	}

	var data []types.DPOSRankLog
	_ = json.Unmarshal(resp.Data, &data)

	hd.responseWrite(ctx, true, data)
}
