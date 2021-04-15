package bean

import "github.com/pkg/errors"

type ContractInvokeTx struct {
	GasLimit        uint64 `json:"gasLimit"` // gas限额
	GasPrice        string `json:"gasPrice"` // gas价格
	Sender          string `json:"sender"`   // 交易发起者公钥
	Privkey         string `json:"privkey"`  // 交易发起者私钥
	Value           string `json:"value"`    // 交易金额
	ContractAddress string `json:"contract"` // 合约地址
	Payload         string `json:"payload"`  // 合约负载 abi.pack(function+参数)
	Memo            string `json:"memo"`     // 备注
}

func (tx *ContractInvokeTx) Check() error {
	if len(tx.Sender) == 0 || len(tx.Privkey) == 0 {
		return errors.New("empty sender or privkey")
	}
	if len(tx.ContractAddress) == 0 {
		return errors.New("empty contract address")
	}
	if tx.GasLimit == 0 {
		return errors.New("ignore gasLimit")
	}
	if tx.GasPrice == "0" {
		return errors.New("ignore gasPrice")
	}
	if len(tx.Memo) > 255 {
		return errors.New("memo length too long")
	}
	return nil
}
