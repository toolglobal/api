package bean

import "github.com/pkg/errors"

type ContractCallTx struct {
	GasLimit        uint64 `json:"gasLimit"` // gas限额
	GasPrice        string `json:"gasPrice"` // gas价格
	Sender          string `json:"sender"`   // 交易发起者公钥
	Privkey         string `json:"privkey"`  // 交易发起者私钥
	Value           string `json:"value"`    // 金额，通常为0
	ContractAddress string `json:"contract"` // 合约地址
	Payload         string `json:"payload"`  // 负载数据，abi.pack(function+参数) hex编码字符串
}

func (tx *ContractCallTx) Check() error {
	if len(tx.Sender) == 0 || len(tx.Privkey) == 0 {
		return errors.New("empty sender or privkey")
	}
	if len(tx.ContractAddress) == 0 {
		return errors.New("empty contract address")
	}
	return nil
}
