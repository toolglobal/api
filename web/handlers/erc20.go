package handlers

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
	"github.com/wolot/api/utils"
	"github.com/wolot/api/web/bean"
	"math/big"
	"strings"
)

func (hd *Handler) BalanceOf(ctx *gin.Context) {
	token := ctx.Param("token")
	to := ctx.Param("to")

	tx := types.NewTxEvm()
	tx.CreatedAt = 0
	tx.GasLimit = 100000
	tx.GasPrice = big.NewInt(1)
	tx.Nonce = 1

	privkey, _ := crypto.GenerateKey()
	buff := make([]byte, 32)
	copy(buff[32-len(privkey.D.Bytes()):], privkey.D.Bytes())

	tx.Sender.SetBytes(crypto.PubkeyToAddress(privkey.PublicKey).Bytes())
	tx.Body.To.SetBytes(ethcmn.HexToAddress(token).Bytes())
	tx.Body.Value = big.NewInt(0)
	abiIns, _ := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	tx.Body.Load, _ = abiIns.Pack("balanceOf", ethcmn.HexToAddress(to))
	tx.Signature, _ = tx.Sign(ethcmn.Bytes2Hex(buff))

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
	results, err := abiIns.Unpack("balanceOf", utils.HexToBytes(evmResult.Ret))
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	if len(results) != 1 {
		hd.responseWrite(ctx, false, "Wrong contract execute result")
		return
	}
	balance := results[0].(*big.Int)
	hd.responseWrite(ctx, true, balance.String())
}
