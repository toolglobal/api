package types

var (
	TxTagAppInit        = []byte{0, 0}     // 账户迁移，采用与batch相同的交易结构，但是不收取手续费 v4废弃
	TxTagTinInit        = []byte{0, 1}     // v3初始化TIN用户 v4废弃
	TxTagAppOLO         = []byte{1, 1}     // 原生交易 v1 - v3废弃
	TxTagAppEvm         = []byte{1, 2}     // 合约交易
	TxTagAppFee         = []byte{1, 3}     // 收取手续费 - v3废弃
	TxTagAppBatch       = []byte{1, 4}     // 批量交易
	TxTagNodeDelegate   = []byte{2, 1}     // 节点抵押赎回提现
	TxTagUserDelegate   = []byte{2, 2}     // 用户抵押赎回提现
	TxTagAppEvmMultisig = []byte{3, 2}     // EVM多签交易
	TxTagAppParams      = []byte{255, 0}   // 修改APP参数
	TxTagAppMgr         = []byte{255, 255} // 节点管理交易
)
