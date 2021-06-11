package handlers

import (
	"encoding/hex"
	"encoding/json"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
	"github.com/wolot/api/utils"
	"github.com/wolot/api/web/bean"
	"go.uber.org/zap"
	"math/big"
	"time"
)

// @Summary 发起evm交易
// @Description 发起evm交易，支持普通转账、合约交易，普通转账默认使用本接口
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.SignedEvmTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/contract/transactions [post]
func (hd *Handler) SignedEvmTransaction(ctx *gin.Context) {
	var tdata bean.SignedEvmTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.logger.Debug("SignedEvmTransaction", zap.Any("tdata", tdata))
	if err := tdata.Check(); err != nil {
		hd.logger.Warn("Check", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	//b, _ := json.Marshal(&tdata)
	//fmt.Println(string(b))
	tx := types.NewTxEvm()
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
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}

	if len(tdata.Body.To) > 42 {
		pub, err := types.HexToPubkey(tdata.Body.To)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(tdata.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(tdata.Body.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Body.Load)
	tx.Body.Memo = []byte(tdata.Body.Memo)
	tx.Signature = utils.HexToBytes(tdata.Signature)

	if !tx.Verify() {
		hd.logger.Warn("Verify failed", zap.Any("tdata", tdata))
		hd.responseWrite(ctx, false, "API SignCheck Failed")
		return
	}

	switch tdata.Mode {
	case bean.MODE_ASYNC:
		hd.signedEvmSendToAsyncTx(ctx, types.TxTagAppEvm[:], tx)
	case bean.MODE_SYNC:
		hd.signedEvmSendToSyncTx(ctx, types.TxTagAppEvm[:], tx)
	default:
		hd.signedEvmSendToCommitTx(ctx, types.TxTagAppEvm[:], tx)
	}
}

func (hd *Handler) signedEvmSendToCommitTx(ctx *gin.Context, tag []byte, sigTx ITx) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()

	result, err := hd.client.BroadcastTxCommit(ctx, append(tag, txBytes...))
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
	response["logs"] = result.DeliverTx.Info
	parseResult(response, result.DeliverTx.Data)

	if result.DeliverTx.Code != types.CodeType_OK {
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.CheckTx.Code))
		hd.responseWriteV2(ctx, false, response, result.DeliverTx.Log)
		return
	}

	hd.responseWriteV2(ctx, true, response, "")
}

func (hd *Handler) signedEvmSendToSyncTx(ctx *gin.Context, tag []byte, sigTx ITx) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()
	result, err := hd.client.BroadcastTxSync(ctx, append(tag, txBytes...))
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

func (hd *Handler) signedEvmSendToAsyncTx(ctx *gin.Context, tag []byte, sigTx ITx) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()

	result, err := hd.client.BroadcastTxAsync(ctx, append(tag, txBytes...))
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

// @Summary 部署合约（非离线签名，慎用）
// @Description 部署合约（私钥交易，非离线签名，慎用）
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.ContractDeployTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/contract/deploy [post]
func (hd *Handler) ContractDeployTx(ctx *gin.Context) {
	var tdata bean.ContractDeployTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	var err error
	tx := types.NewTxEvm()
	tx.CreatedAt = uint64(time.Now().UnixNano())
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	if len(tdata.Sender) > 42 {
		pub, err := types.HexToPubkey(tdata.Sender)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}
	act, err := hd.v2QueryAccount(tx.Sender.ToAddress().Address.Hex())
	if err != nil {
		hd.responseWrite(ctx, false, "get nonce:"+err.Error())
		return
	}
	tx.Nonce = act.Nonce
	tx.Body.Value, _ = new(big.Int).SetString(tdata.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Payload)
	tx.Body.Memo = []byte(tdata.Memo)
	tx.Signature, err = tx.Sign(tdata.Privkey)
	if err != nil {
		hd.responseWrite(ctx, false, "tx.Sign:"+err.Error())
		return
	}
	hd.contractDeployCommitTx(ctx, tx)
}

func (hd *Handler) contractDeployCommitTx(ctx *gin.Context, sigTx *types.TxEvm) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()
	response["address"] = crypto.CreateAddress(sigTx.Sender.ToAddress().Address, sigTx.Nonce).Hex()

	result, err := hd.client.BroadcastTxCommit(ctx, append(types.TxTagAppEvm[:], txBytes...))
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
		hd.logger.Info("DeliverTx", zap.Uint32("code", result.DeliverTx.Code))
		hd.responseWriteV2(ctx, false, response, result.DeliverTx.Log)
		return
	}

	hd.responseWriteV2(ctx, true, response, "")
}

// @Summary 调用合约（非离线签名，慎用）
// @Description 调用合约（私钥交易，非离线签名，慎用）,使用本接口调用合约只读方法会消耗gas。
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.ContractInvokeTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/contract/invoke [post]
func (hd *Handler) ContractInvokeTx(ctx *gin.Context) {
	var tdata bean.ContractInvokeTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	var err error
	tx := types.NewTxEvm()
	tx.CreatedAt = uint64(time.Now().UnixNano())
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	if len(tdata.Sender) > 42 {
		pub, err := types.HexToPubkey(tdata.Sender)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}
	act, err := hd.v2QueryAccount(tx.Sender.ToAddress().Address.Hex())
	if err != nil {
		hd.responseWrite(ctx, false, "get nonce:"+err.Error())
		return
	}
	tx.Nonce = act.Nonce
	if len(tdata.ContractAddress) > 42 {
		pub, err := types.HexToPubkey(tdata.ContractAddress)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(tdata.ContractAddress).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(tdata.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Payload)
	tx.Body.Memo = []byte(tdata.Memo)
	tx.Signature, err = tx.Sign(tdata.Privkey)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	hd.contractInvokeCommitTx(ctx, tx)
}

// @Summary 只读调用合约（非离线签名，慎用）
// @Description 仅在本节节点evm副本上执行合约，不会广播交易，不对区块产生影响，不消耗gas。用户估算智能合约消耗的gas，以及用于调用合约只读方法。
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.ContractCallTx true "请求参数"
// @Success 200 {object} bean.EvmCallResult "成功"
// @Router /v2/contract/call [post]
func (hd *Handler) ContractCallTx(ctx *gin.Context) {
	var tdata bean.ContractCallTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	var err error
	tx := types.NewTxEvm()
	tx.CreatedAt = uint64(time.Now().UnixNano())
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	if len(tdata.Sender) > 42 {
		pub, err := types.HexToPubkey(tdata.Sender)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Sender.SetBytes(pub.Bytes())
	} else {
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}
	act, err := hd.v2QueryAccount(tx.Sender.ToAddress().Address.Hex())
	if err != nil {
		hd.responseWrite(ctx, false, "get nonce:"+err.Error())
		return
	}
	tx.Nonce = act.Nonce
	if len(tdata.ContractAddress) > 42 {
		pub, err := types.HexToPubkey(tdata.ContractAddress)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(tdata.ContractAddress).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(tdata.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Payload)

	tx.Signature, err = tx.Sign(tdata.Privkey)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	hd.callContract(ctx, tx)
}

// @Summary 只读调用合约
// @Description 仅在本节节点evm副本上执行合约，不会广播交易，不对区块产生影响，不消耗gas。用户估算智能合约消耗的gas，以及用于调用合约只读方法。
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.SignedEvmTx true "请求参数"
// @Success 200 {object} bean.EvmCallResult "成功"
// @Router /v2/contract/query [post]
func (hd *Handler) ContractSignedCallTx(ctx *gin.Context) {
	var tdata bean.SignedEvmTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	tx := types.NewTxEvm()
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
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}

	if len(tdata.Body.To) > 42 {
		pub, err := types.HexToPubkey(tdata.Body.To)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(tdata.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(tdata.Body.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Body.Load)
	tx.Body.Memo = []byte(tdata.Body.Memo)
	tx.Signature = utils.HexToBytes(tdata.Signature)

	if !tx.Verify() {
		hd.responseWrite(ctx, false, "API SignCheck Failed")
		return
	}

	hd.callContract(ctx, tx)
}

func (hd *Handler) callContract(ctx *gin.Context, tx *types.TxEvm) {
	result, err := hd.client.ABCIQuery(ctx, types.API_V2_CONTRACT_CALL, tx.ToBytes())
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
	var evmResult bean.EvmCallResult
	if err = json.Unmarshal(resp.Data, &evmResult); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.responseWrite(ctx, true, evmResult)
}

func (hd *Handler) contractInvokeCommitTx(ctx *gin.Context, sigTx *types.TxEvm) {
	var (
		txBytes  = sigTx.ToBytes()
		response = make(map[string]interface{})
	)
	response["tx"] = sigTx.Hash().Hex()

	result, err := hd.client.BroadcastTxCommit(ctx, append(types.TxTagAppEvm[:], txBytes...))
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

func parseResult(response map[string]interface{}, bs []byte) {
	var result types.TxResult
	if err := result.FromBytes(bs); err == nil {
		response["ret"] = hex.EncodeToString(result.Ret)
		//response["logs"] = result.Logs
	}
}

// @Summary 根据transaction参数进行签名
// @Description 根据transaction参数进行签名
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.SignEvmTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/contract/sign [post]
func (hd *Handler) SignEvmTransaction(ctx *gin.Context) {
	var tdata bean.SignEvmTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	if err := tdata.Check(); err != nil {
		hd.logger.Warn("Check", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	tx := types.NewTxEvm()
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
		tx.Sender.SetBytes(ethcmn.HexToAddress(tdata.Sender).Bytes())
	}

	if len(tdata.Body.To) > 42 {
		pub, err := types.HexToPubkey(tdata.Body.To)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		tx.Body.To.SetBytes(pub.Bytes())
	} else {
		tx.Body.To.SetBytes(ethcmn.HexToAddress(tdata.Body.To).Bytes())
	}

	tx.Body.Value, _ = new(big.Int).SetString(tdata.Body.Value, 10)
	tx.Body.Load = utils.HexToBytes(tdata.Body.Load)
	tx.Body.Memo = []byte(tdata.Body.Memo)
	sign, err := tx.Sign(tdata.PrivateKey)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.responseWrite(ctx, true, hex.EncodeToString(sign))
}

// @Summary 发起evm多签交易
// @Description 发起evm多签交易，支持OLO转账、合约交易
// @Tags v2-contract
// @Accept json
// @Produce json
// @Param Request body bean.SignedMultisigEvmTx true "请求参数"
// @Success 200 {object}  bean.PublicResp "成功"
// @Router /v2/contract/multisigTransactions [post]
func (hd *Handler) SignedEvmMutlisigTransaction(ctx *gin.Context) {
	var tdata bean.SignedMultisigEvmTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.logger.Debug("SignedEvmMutlisigTransaction", zap.Any("tdata", tdata))
	if err := tdata.Check(); err != nil {
		hd.logger.Warn("Check", zap.Error(err))
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var pkeys []types.PublicKey
	for _, v := range tdata.Signature.PubKey.PubKeys {
		var pkey types.PublicKey
		copy(pkey[:], utils.HexToBytes(v))
		pkeys = append(pkeys, pkey)
	}

	tx := types.NewMultisigEvmTx(tdata.Signature.PubKey.K, pkeys)
	tx.Deadline = tdata.Deadline
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	tx.From = ethcmn.HexToAddress(tdata.From)
	tx.Nonce = tdata.Nonce
	tx.To = ethcmn.HexToAddress(tdata.To)
	tx.Value, _ = new(big.Int).SetString(tdata.Value, 10)
	tx.Load = utils.HexToBytes(tdata.Load)
	tx.Memo = []byte(tdata.Memo)

	for _, v := range tdata.Signature.Signatures {
		if err := tx.AddSign(v); err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
	}

	if !tx.Verify() {
		hd.logger.Warn("Verify failed", zap.Any("tdata", tdata))
		hd.responseWrite(ctx, false, "API SignCheck Failed")
		return
	}

	switch tdata.Mode {
	case bean.MODE_ASYNC:
		hd.signedEvmSendToAsyncTx(ctx, types.TxTagAppEvmMultisig[:], tx)
	case bean.MODE_SYNC:
		hd.signedEvmSendToSyncTx(ctx, types.TxTagAppEvmMultisig[:], tx)
	default:
		hd.signedEvmSendToCommitTx(ctx, types.TxTagAppEvmMultisig[:], tx)
	}
}

func (hd *Handler) Multisigner(ctx *gin.Context) {
	var tdata bean.SignedMultisigEvmTx
	if err := ctx.ShouldBindJSON(&tdata); err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	hd.logger.Debug("Multisigner", zap.Any("tdata", tdata))

	var pkeys []types.PublicKey
	for _, v := range tdata.Signature.PubKey.PubKeys {
		var pkey types.PublicKey
		copy(pkey[:], utils.HexToBytes(v))
		pkeys = append(pkeys, pkey)
	}

	tx := types.NewMultisigEvmTx(tdata.Signature.PubKey.K, pkeys)
	tx.Deadline = tdata.Deadline
	tx.GasLimit = tdata.GasLimit
	tx.GasPrice, _ = new(big.Int).SetString(tdata.GasPrice, 10)
	tx.From = ethcmn.HexToAddress(tdata.From)
	tx.Nonce = tdata.Nonce
	tx.To = ethcmn.HexToAddress(tdata.To)
	tx.Value, _ = new(big.Int).SetString(tdata.Value, 10)
	tx.Load = utils.HexToBytes(tdata.Load)
	tx.Memo = []byte(tdata.Memo)

	for _, v := range tdata.Signature.Signatures {
		if err := tx.AddSign(v); err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
	}

	var signers []string
	for _, v := range tx.Signer() {
		signers = append(signers, v.Hex())
	}
	hd.responseWrite(ctx, true, signers)
}

type ITx interface {
	ToBytes() []byte
	Hash() ethcmn.Hash
}
