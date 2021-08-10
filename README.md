# API
本API程序轮询区块，解析并生成可读的block、transaction、payments数据存入sqlite，提供http REST接口方便查询。

由于新的mondod提供兼容以太坊web3 JSON RPC，本API不再长期支持，请使用者优先使用JSON RPC。

## Install
```shell
git clone https://github.com/toolglobal/api.git
cd api && make
cd build && ./api
```

## config
```toml
bind = ":8889" # 监听8889 http端口
rpc = "127.0.0.1:26657" # 连接本地mondod节点的26657 tendermint rpc端口
dev = true # 开发模式
metrics = true # prometheus 监控
chainId = "8723" # 链id，mainnet：8723 testnet：8724
versions = [3] # 解析协议版本
startHeight = 1 # 开始解析区块高度
tgsBaseURL = "https://services.wolot.io" # 获取官方代币配置的接口

[limiter] # 合约查询限流，合约查询需要执行evm，性能损耗大，可能影响节点稳定
interval = "0h0m1s"
capacity = 100
```