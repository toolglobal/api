package types

import (
	"errors"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

const (
	CodeType_OK uint32 = 0
	// General response codes, 0 ~ 99
	CodeType_InternalError uint32 = 1 // 内部错误
	CodeType_EncodingError uint32 = 2 // 编解码错误
	CodeType_UnsupportedTx uint32 = 3 // 不支持的tx类型
	CodeType_BadSignature  uint32 = 4 // 签名错误
	CodeType_BadArgument   uint32 = 5 // 参数错误
	CodeType_ExecuteTx     uint32 = 6 // tx执行错误
	CodeType_Deadline      uint32 = 7 // 交易过期
	CodeType_Gas           uint32 = 8 // gas太低

	CodeType_SenderNotExist    uint32 = 11 // 帐户不存在
	CodeType_BadNonce          uint32 = 12 // nonce错误
	CodeType_InsufficientFunds uint32 = 13 // 资金不足
	CodeType_ReceiverNotExist  uint32 = 14 // 帐户不存在

	// evm tx
	CodeType_ContractExecuteErr    uint32 = 30 // 合约执行错误
	CodeType_ContractExecuteFailed uint32 = 31 // 合约执行失败
)

const (
	TxGas = int64(21000) // tx gas
)

var (
	ZERO_ADDRESS   = ethcmn.Address{}
	ZERO_HASH      = ethcmn.Hash{}
	ZERO_PUBLICKEY = PublicKey{}
)

var (
	ErrUnsupportedTx           = errors.New("unsupported tx")
	ErrInvalidSignature        = errors.New("signature verify failed")
	ErrInvalidOperations       = errors.New("too little/much operations")
	ErrInvalidGasLimit         = errors.New("invalid gasLimit")
	ErrInvalidGasPrice         = errors.New("invalid gasPrice")
	ErrNotFoundSender          = errors.New("sender not exist")
	ErrNotFoundReceiver        = errors.New("receiver not exist")
	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrDeadline                = errors.New("tx deadline")
	ErrGas                     = errors.New("gas too low")
)
