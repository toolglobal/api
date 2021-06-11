package handlers

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
	"github.com/wolot/api/web/bean"
	"go.uber.org/zap"
	"math/big"
)

// @Summary 用户委托交易
// @Description 用户委托交易
// @Tags v3-dpos
// @Accept json
// @Produce json
// @Param Request body bean.DelegateTx true "请求参数"
// @Success 200 "成功"
// @Router /v3/delegate [post]
func (hd *Handler) V3_SendDelegateTx(ctx *gin.Context) {
	var tdata bean.DelegateTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	tx := types.NewUserDelegateTx()
	tx.CreatedAt = tdata.CreatedAt
	if len(tdata.Sender) > 42 {
		pub, err := types.HexToPubkey(tdata.Sender)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(common.HexToAddress(tdata.Sender).Bytes())
	}
	tx.Nonce = tdata.Nonce
	tx.OpType = tdata.OpType
	tx.OpValue, _ = new(big.Int).SetString(tdata.OpValue, 10)
	if tdata.Receiver != "" {
		tx.Receiver = common.HexToAddress(tdata.Receiver).Bytes()
	}
	tx.Signature, _ = hex.DecodeString(tdata.Signature)

	if !tx.Verify() {
		hd.responseWrite(ctx, false, "API SignCheck Failed")
		return
	}
	hd.signedDelegateSendToCommitTx(ctx, tx)
}

func (hd *Handler) signedDelegateSendToCommitTx(ctx *gin.Context, sigTx *types.TxUserDelegate) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()

	result, err := hd.client.BroadcastTxCommit(ctx, append(types.TxTagUserDelegate[:], txBytes...))
	if err != nil {
		hd.logger.Error("BroadcastTxCommit", zap.Error(err))
		hd.responseWriteV2(ctx, false, response, err.Error())
		return
	}
	if result.CheckTx.Code != types.CodeType_OK {
		hd.logger.Info("CheckTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWriteV2(ctx, false, response, result.CheckTx.Log)
		return
	}

	if result.DeliverTx.Code != types.CodeType_OK {
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWriteV2(ctx, false, response, result.DeliverTx.Log)
		return
	}

	hd.responseWriteV2(ctx, true, response, "")
}
