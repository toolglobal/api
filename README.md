# API

## Install
```shell
git clone https://github.com/toolglobal/api.git
cd api && make
cd build && ./api
```
## config
```toml
bind = ":8889"
rpc = "127.0.0.1:26657" # mondo tendermint rpc address
dev = true # dev mode,open api docs http://$bind$/docs/index.html
metrics = true # gin prometheus
startHeight = 1 # sync block from startHeight
versions = [3] # only version 3 is valid

[limiter] # 可选 rate limiter,for /v2/contract/query
interval = "0h0m1s"
capacity = 10

[[tokens.coins]] # 可选，如果配置了token信息，API程序将根据代币合约地址解析该token的转账信息写入payment表，仅支持ERC20 token
name = "Bridge Digital Economy"
symbol = "BDE"
address = "0x67EbBA731DCd5b763F5699650920a637eDbBEb93"
decimals = 8
icon = "http://127.0.0.1:8889/static/coin/BDE.ico"
```