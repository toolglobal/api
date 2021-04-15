package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/wolot/api/mondo/types"
	"github.com/wolot/api/web/bean"
)

// @Summary 随机生成mondo账户
// @Description 随机生成mondo账户，该账户默认未上链，需要通过转账交易激活
// @Tags v2-general
// @Accept json
// @Produce json
// @Success 200 {object}  bean.V2GenKeyResult "成功"
// @Router /v2/genkey [get]
func (hd *Handler) V2GenKey(ctx *gin.Context) {
	privkey, err := crypto.GenerateKey()
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	var result bean.V2GenKeyResult

	buff := make([]byte, 32)
	copy(buff[32-len(privkey.D.Bytes()):], privkey.D.Bytes())
	result.Privkey = ethcmn.Bytes2Hex(buff)
	result.Pubkey = ethcmn.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey))
	result.Address = crypto.PubkeyToAddress(privkey.PublicKey).String()

	hd.responseWrite(ctx, true, result)
}

// @Summary 查询账户信息
// @Description 根据地址查询账户信息
// @Tags v2-general
// @Accept json
// @Produce json
// @Param address path string true "账户地址"
// @Success 200 {object}  bean.V2AccountResult "成功"
// @Router /v2/accounts/{address} [get]
func (hd *Handler) V2QueryAccount(ctx *gin.Context) {
	addressHex := ctx.Param("address")

	if len(addressHex) > 42 {
		pubkey, err := types.HexToPubkey(addressHex)
		if err != nil {
			hd.responseWrite(ctx, false, err.Error())
			return
		}
		addressHex = pubkey.ToAddress().Address.Hex()
	}

	act, err := hd.v2QueryAccount(addressHex)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}

	hd.responseWrite(ctx, true, act)
}

func (hd *Handler) v2QueryAccount(address string) (*bean.V2AccountResult, error) {
	result, err := hd.client.ABCIQuery(context.Background(), types.API_V2_QUERY_ACCOUNT, ethcmn.HexToAddress(address).Bytes())
	if err != nil {
		return nil, err
	}

	var data types.Result
	err = rlp.DecodeBytes(result.Response.Value, &data)
	if err != nil {
		return nil, err
	}
	if data.Code != types.CodeType_OK {
		return nil, fmt.Errorf("code %d, log %s", data.Code, data.Log)
	}
	var act bean.V2AccountResult
	if err := json.Unmarshal(data.Data, &act); err != nil {
		return nil, err
	}
	return &act, nil
}

// @Summary 地址转换
// @Description 地址转换，公钥转地址
// @Tags v2-general
// @Accept json
// @Produce json
// @Param publickey path string true "v1账户地址，即用户公钥"
// @Success 200 {object}  bean.V2ConvertResult "成功"
// @Router /v2/convert/{publickey} [get]
func (hd *Handler) V2Convert(ctx *gin.Context) {
	publickey := ctx.Param("publickey")
	if len(publickey) != 66 && len(publickey) != 68 {
		hd.responseWrite(ctx, false, "error length of public key")
		return
	}
	pubkey, err := types.HexToPubkey(publickey)
	if err != nil {
		hd.responseWrite(ctx, false, err.Error())
		return
	}
	pubkey.ToAddress()
	resData := bean.V2ConvertResult{
		OldAddress: publickey,
		PublicKey:  publickey,
		Address:    types.PubkeyToAddress(pubkey).String(),
	}
	hd.responseWrite(ctx, true, &resData)
}
