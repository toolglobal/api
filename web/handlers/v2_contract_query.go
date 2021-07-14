package handlers

import (
	"encoding/json"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/toolglobal/api/mondo/types"
	"github.com/toolglobal/api/web/bean"
)

// @Summary 查询合约账户信息
// @Description 查询合约账户信息
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param address path string true "合约地址"
// @Success 200 {object}  bean.V2ContractActResult "成功"
// @Router /v2/contract/accounts/{address} [get]
func (hd *Handler) QueryContract(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if len(addressHex) > 42 {
		pubkey, err := types.HexToPubkey(addressHex)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		addressHex = pubkey.ToAddress().Address.Hex()
	}

	result, err := hd.client.ABCIQuery(ctx, types.API_V2_CONTRACT_QUERY_ACCOUNT, ethcmn.HexToAddress(addressHex).Bytes())
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

	show := bean.V2ContractActResult{}
	json.Unmarshal(resp.Data, &show)

	hd.responseWrite(ctx, true, show)
}

// @Summary 查询事件日志
// @Description 查询事件日志
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param txhash path string true "交易hash"
// @Success 200 {object}  bean.V2EvmLog "成功"
// @Router /v2/contract/events/{txhash} [get]
func (hd *Handler) QueryTxEvents(ctx *gin.Context) {
	hash := ctx.Param("txhash")

	result, err := hd.client.ABCIQuery(ctx, types.API_V2_CONTRACT_QUERY_LOGS, ethcmn.HexToHash(hash).Bytes())
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

	hd.responseWrite(ctx, true, string(resp.Data))
}
