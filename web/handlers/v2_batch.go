package handlers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/toolglobal/api/mondo/types"
	"github.com/toolglobal/api/utils"
	"github.com/toolglobal/api/web/bean"
	"go.uber.org/zap"
	"math/big"
)

// @Summary 发起批量交易
// @Description 发起批量交易
// @Tags v2-general
// @Accept json
// @Produce json
// @Param Request body bean.SignedBatchTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/transactions [post]
func (hd *Handler) SignedBatchTransaction(ctx *gin.Context) {
	var tdata bean.SignedBatchTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	tx := types.NewTxBatch()
	tx.CreatedAt = tdata.CreatedAt
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	tx.Nonce = tdata.Nonce

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

	tx.Memo = []byte(tdata.Memo)
	for _, v := range tdata.Ops {
		var to types.PublicKey
		if len(v.To) > 42 {
			pub, err := types.HexToPubkey(v.To)
			if err != nil {
				hd.responseWrite(ctx, false, err.Error())
				return
			}
			to.SetBytes(pub.Bytes())
		} else {
			to.SetBytes(common.HexToAddress(v.To).Bytes())
		}
		value, _ := new(big.Int).SetString(v.Value, 10)
		op := types.TxOp{
			To:    to,
			Value: value,
		}
		tx.Ops = append(tx.Ops, op)
	}
	tx.Signature = utils.HexToBytes(tdata.Signature)

	if !tx.Verify() {
		hd.responseWrite(ctx, false, "API SignCheck Failed")
		return
	}
	switch tdata.Mode {
	case bean.MODE_ASYNC:
		hd.signedBatchSendToAsyncTx(ctx, tx)
	case bean.MODE_SYNC:
		hd.signedBatchSendToSyncTx(ctx, tx)
	default:
		hd.signedBatchSendToCommitTx(ctx, tx)
	}
}

func (hd *Handler) signedBatchSendToCommitTx(ctx *gin.Context, sigTx *types.TxBatch) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()
	result, err := hd.client.BroadcastTxCommit(ctx, append(types.TxTagAppBatch[:], txBytes...))
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
	response["code"] = result.DeliverTx.Code
	response["gasUsed"] = result.DeliverTx.GasUsed
	parseResult(response, result.DeliverTx.Data)

	if result.DeliverTx.Code != types.CodeType_OK {
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWriteV2(ctx, false, response, result.DeliverTx.Log)
		return
	}

	hd.responseWriteV2(ctx, true, response, "")
}

func (hd *Handler) signedBatchSendToSyncTx(ctx *gin.Context, sigTx *types.TxBatch) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()
	result, err := hd.client.BroadcastTxSync(ctx, append(types.TxTagAppBatch[:], txBytes...))
	if err != nil {
		hd.logger.Error("BroadcastTxSync", zap.Error(err))
		hd.responseWriteV2(ctx, false, response, err.Error())
		return
	}
	if result.Code != types.CodeType_OK {
		hd.logger.Info("BroadcastTxSync", zap.Uint32("code", result.Code), zap.String("log", result.Log))
		hd.responseWriteV2(ctx, false, response, result.Log)
		return
	}

	hd.responseWriteV2(ctx, true, response, "")
}

func (hd *Handler) signedBatchSendToAsyncTx(ctx *gin.Context, sigTx *types.TxBatch) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()
	result, err := hd.client.BroadcastTxAsync(ctx, append(types.TxTagAppBatch[:], txBytes...))
	if err != nil {
		hd.logger.Error("BroadcastTxAsync", zap.Error(err))
		hd.responseWriteV2(ctx, false, response, err.Error())
		return
	}
	if result.Code != types.CodeType_OK {
		hd.logger.Info("BroadcastTxAsync", zap.Uint32("code", result.Code), zap.String("log", result.Log))
		hd.responseWriteV2(ctx, false, response, result.Log)
		return
	}
	hd.responseWriteV2(ctx, true, response, "")
}
