package handlers

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
)

// @Summary 查询投票者信息
// @Description 查询投票者信息
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param address path string true "地址"
// @Success 200 {object}  types.BeanVoter "成功"
// @Router /v3/voters/{address} [get]
func (hd *Handler) V3QueryVoter(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if len(addressHex) != 42 {
		hd.responseWrite(ctx, false, "invalid address")
		return
	}

	result, err := hd.client.ABCIQuery(ctx, types.API_V3_QUERY_VOTER, common.HexToAddress(addressHex).Bytes())
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

	var voter types.BeanVoter
	_ = json.Unmarshal(resp.Data, &voter)

	hd.responseWrite(ctx, true, voter)
}
