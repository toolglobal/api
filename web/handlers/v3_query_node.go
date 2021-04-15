package handlers

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
)

// @Summary 查询节点信息
// @Description 查询节点信息
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address path string true "地址"
// @Success 200 {object}  types.BeanValidator "成功"
// @Router /v3/nodes/{address} [get]
func (hd *Handler) V3QueryNode(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if len(addressHex) != 40 {
		hd.responseWrite(ctx, false, "invalid address")
		return
	}

	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_NODE, common.HexToAddress(addressHex).Bytes())
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

	var node types.BeanValidator
	_ = json.Unmarshal(resp.Data, &node)

	hd.responseWrite(ctx, true, node)
}

// @Summary 查询节点投票人信息
// @Description 查询节点投票人信息
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address path string true "地址"
// @Success 200 {array}  types.BeanVots "成功"
// @Router /v3/nodes/{address}/voters [get]
func (hd *Handler) V3QueryNodeVoters(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if len(addressHex) != 40 {
		hd.responseWrite(ctx, false, "invalid address")
		return
	}

	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_NODEVOTERS, common.HexToAddress(addressHex).Bytes())
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

	var node types.BeanVots
	_ = json.Unmarshal(resp.Data, &node)

	hd.responseWrite(ctx, true, node)
}

// @Summary 查询所有节点信息
// @Description 查询所有节点信息
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Success 200 {array}  types.BeanValidator "成功"
// @Router /v3/nodes [get]
func (hd *Handler) V3QueryNodes(ctx *gin.Context) {
	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_NODES, nil)
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

	var nodes []types.BeanValidator
	_ = json.Unmarshal(resp.Data, &nodes)

	hd.responseWrite(ctx, true, nodes)
}
