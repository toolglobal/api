# MondoAPI_V3

[TOC]
# 1. 概述
本文档时mondo 链的api程序对外提供的api文档，api程序已开源用户可自行修改。api程序本身提供swag格式的文档输出（https://oloapi-test.wolot.io/docs/）。

调用链路：【用户】----http----【API】----RPC----【OLO CHAIN】
如上所示，本API是在用户与区块链之间的桥梁，使得用户可以通过简单的http请求链的RPC接口，比如账户查询、交易等；同时，API程序沦陷OLO CHAIN，生成"leadger","transaction","payment"三种数据，提供给用户查询。

## 1.1 公共应答结构
```
type Response struct {
	IsSuccess bool        `json:"isSuccess"`//是否成功
	Result    interface{} `json:"result"`   //返回的数据
	Message   string      `json:"message"`  //消息提示
}
```

## 1.2 查询公共控制参数（query）
字段       |字段类型       |字段说明
------------|-----------|-----------
cursor       |int           |游标，从0开始
limit          |int              |每页条数（默认10，最大200）
order        |string        |排序方式(asc/desc，默认desc)
begin        |int              |开始时间戳
end           |int              |结束时间戳
height       |int64          |区块高度
symbol      |string          |币种，例如"OLO","ABC"
address     |string          |用户地址或节点地址

## 1.3 交易签名
OLO链的私钥算法与以太坊一致，地址、私钥与以太坊互通，签名也采用类似的方法，区别在于OLO链的交易结构不同。

### 签名步骤
- 构造交易对象 TxEvm
- 对交易参数进行RLP编码
```
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.GasPrice,
		tx.GasLimit,
		tx.Nonce,
		tx.Sender,
		tx.Body,
	})
```
- 对上一步RLP编码结果计算RLP哈希，hash256
```
func rlpHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

此方法可以拆解为两步：
1. 对上一步结果进一步进行RLP编码，得到编码后字节数组
2. 对此字节数组计算hash256
```
- 用私钥对上一步结果hash值进行签名
```
crypto.Sign(tx.SigHash().Bytes(), priv)
```
- 将签名填充到交易对象的Signature中

### 计算离线hash
- 构造交易对象 TxEvm
- 对整个交易结构所有参数计算RLP HASH，即为交易hash
```
v := rlpHash(tx)
func rlpHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
```

### Golang版签名源码
```
type TxEvm struct {
	CreatedAt uint64      // 交易发起时间
	GasLimit  uint64      // gas限额
	GasPrice  *big.Int    // gas价格
	Nonce     uint64      // 交易发起者nonce
	Sender    PublicKey   // 交易发起者公钥 兼容v1
	Body      TxEvmCommon // 交易结构
	Signature []byte      // 交易签名
	hash      atomic.Value
}

const PubKeyLength = 33
type PublicKey [PubKeyLength]byte

type TxEvmCommon struct {
	To    PublicKey // 交易接收方公钥或地址，地址时填后20字节，创建合约是为全0；兼容v1
	Value *big.Int  // 交易金额
	Load  []byte    // 合约交易负荷
	Memo  []byte    // 备注信息
}

func NewTxEvm() *TxEvm {
	return &TxEvm{}
}

func (tx *TxEvm) Sign(privkey string) ([]byte, error) {
	if strings.HasPrefix(privkey, "0x") {
		privkey = string([]byte(privkey)[2:])
	}

	priv, err := crypto.ToECDSA(ethcmn.Hex2Bytes(privkey))
	if err != nil {
		return nil, err
	}
	return crypto.Sign(tx.SigHash().Bytes(), priv)
}

func (tx *TxEvm) SigHash() ethcmn.Hash {
	message, _ := rlp.EncodeToBytes([]interface{}{
		tx.GasPrice,
		tx.GasLimit,
		tx.Nonce,
		tx.Sender,
		tx.Body,
	})

	return rlpHash(message)
}

func (tx *TxEvm) Hash() ethcmn.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(ethcmn.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

func rlpHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

```

### Java版签名源码
```
import lombok.Data;
import org.ethereum.crypto.ECKey;
import org.ethereum.crypto.HashUtil;
import org.ethereum.util.RLP;
import org.spongycastle.util.encoders.Hex;
import java.math.BigInteger;


@Data
public class Transaction {
    private String gasLimit;
    private String gasPrice;
    private String nonce;
    private String sender;
    private TxBody body;
    private byte[] signature;
    private String privkey;

    private static final int OFFSET_SHORT_LIST = 0xc0;

    // 设置签名
    public void Sign(){
        byte[] signHash = toSignHash();
        sign(signHash);
    }

    private void sign(byte[] hash) {
        ECKey ecKey = new ECKey().fromPrivate(Hex.decode(privkey));
        signature = ecKey.sign(hash).toByteArray();
    }

    private byte[] toSignHash() {
        byte[] rlpRaw = RLPEncode();
        return HashUtil.sha3(RLP.encode(rlpRaw));
    }

    private byte[] RLPEncode(){
        byte[] bodybytes = RLP.encodeList(RLP.encodeElement(Hex.decode(body.getTo())),
                RLP.encodeBigInteger(new BigInteger(body.getValue())),
                RLP.encodeElement(Hex.decode(body.getLoad())),
                RLP.encodeString(body.getMemo())
        );

        return RLP.encodeList(RLP.encode(new BigInteger(gasPrice)),
                RLP.encode(new BigInteger(gasLimit)),
                RLP.encode(new BigInteger(nonce)),
                RLP.encodeElement(Hex.decode(sender)),
                bodybytes
        );
    }

    // 返回交易hash
    public byte[] Hash(){
        byte[] rlpRaw = RLPEncodeTx();
        return HashUtil.sha3(rlpRaw);
    }

    private byte[] RLPEncodeTx(){
        byte[] bodybytes = RLP.encodeList(RLP.encodeElement(Hex.decode(body.getTo())),
                RLP.encodeBigInteger(new BigInteger(body.getValue())),
                RLP.encodeElement(Hex.decode(body.getLoad())),
                RLP.encodeString(body.getMemo())
        );

        return RLP.encodeList(RLP.encode(new BigInteger(gasLimit)),
                RLP.encode(new BigInteger(gasPrice)),
                RLP.encode(new BigInteger(nonce)),
                RLP.encodeElement(Hex.decode(sender)),
                bodybytes,
                RLP.encode(signature)
        );
    }
}
// -------------
@Data
public class TxBody {
    private String to;
    private String value;
    private String load;//hex
    private String memo;
}

```

# 2. v2普通接口(/v2)

## 2.1 生成帐户
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v2/genkey](#)
> 说明：通过固定算法生成帐户，该账户未上链；若要上链，可通过转账方法向本账户转账，会自动创建本账户，转账金额不能为0。此接口可能不安全，建议本地调用。

### 请求参数
无

### 返回结果
```
var result struct {
    Privkey string `json:"privkey"`//私钥
    Pubkey  string `json:"pubkey"` //公钥(压缩公钥)
    Address string `json:"address"`//地址
}
```

### 示例
> http://192.168.8.145:10000/v2/genkey
```
{
	"isSuccess": true,
	"result": {
		"privkey": "fa49a5ddb3693c47255ba8a53bac77f8f40d7c668a84a9a42d384de7934c6e6d",
		"pubkey": "02fa04a9ec89831a7bb5a62f89cad0ffb7b3e9bee9ca3df88a84b0e5ceb42677d9",
		"address": "0x4dDE4aD8b8eFe6446f0A33238F7951b0f71Efc1E"
	}
}
```

## 2.2 根据地址查询帐户信息
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v2/accounts/:address](#)

### 请求参数
字段       |字段类型       |字段说明
------------|-----------|-----------
address       |string        |帐户地址，或公钥

### 返回结果
```
type V2AccountResult struct {
	Address string `json:"address"` // 地址
	Balance string `json:"balance"` // 余额
	Nonce   uint64 `json:"nonce"`   // nonce
}

```
### 示例
> http://192.168.8.145:10000/v2/accounts/0x7752B42608A0f1943c19FC5802cb027E60B4C911
```
{
	"isSuccess": true,
	"result": {
		"address": "0x7752b42608a0f1943c19fc5802cb027e60b4c911",
		"balance": "999908891070887",
		"nonce": 52
	}
}
```

## 2.3 根据公钥获取地址（v1地址转v2地址）
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v2/convert/:publicKey](#)

### 请求参数
字段       |字段类型       |字段说明
------------|-----------|-----------
publicKey       |string        |公压缩钥（v1旧长地址）

### 返回结果
```
type V2ConvertResult struct {
	PublicKey  string `json:"public_key"`  // 公钥，同上
	Address    string `json:"address"`     // 地址
}
```
### 示例
> http://192.168.8.145:10000/v2/convert/0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5
```
{
	"isSuccess": true,
	"result": {
		"public_key": "0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5",
		"address": "0x7752B42608A0f1943c19FC5802cb027E60B4C911"
	}
}
```


## 2.4 批量交易
### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/transactions](#)
> 批量OLO原生币转载，每次最多可发送10000笔

### 请求参数
```
type SignedBatchTx struct {
	Mode      int         `json:"mode"`       // 模式:0-default/commit 1-async 2-sync
	CreatedAt uint64      `json:"createdAt"`  // 时间戳，可选字段，秒/毫秒均可
	GasLimit  uint64      `json:"gasLimit"`   // gas限额 建议：21000*len(ops)
	GasPrice  string      `json:"gasPrice"`   // gas价格，至少为1
	Nonce     uint64      `json:"nonce"`      // 用户nonce，每次交易前从链上获取，每次交易nonce+1
	Sender    string      `json:"sender"`     // 交易发起者公钥
	Ops       []Operation `json:"operations"` // 交易列表，数量不可大于10000笔
	Memo      string      `json:"memo"`       // 备注，必须<256字节
	Signature string      `json:"signature"`  // 交易签名的hex字符串
}

type Operation struct {
	To    string `json:"to"`    // 交易接受方地址，可以是普通用户地址、合约地址、节点账户地址
	Value string `json:"value"` // 交易金额
}
```
### 返回结果
```
{
	"isSuccess": true,
	"result": {
		"res": "",
		"tx": "0x8deda6b70fb5c15fce17d4956eebfec949cab3e697d8c9083fdba853158037a8" //交易hash
	}
}
```

# 3. v2合约交易接口(/v2/contract)


## 3.1 发起签名合约交易（普通转账调用此接口）
本接口适用场景：部署合约、执行合约方法、OLO原生币转账；注：本接口不适应调用合约只读方法，会消耗GAS；

### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/contract/transactions](#)

### 请求参数
```
type SignedEvmTx struct {
	Mode      int    `json:"mode"`      // 交易模式，可选，默认为0；0-同步模式 1-全异步 2-半异步；如果tx执行时间较长、网络不稳定、出块慢，建议使用半异步模式。
	CreatedAt uint64 `json:"createdAt"` // 时间戳，可选
	GasLimit  uint64 `json:"gasLimit"`  // gas限额
	GasPrice  string `json:"gasPrice"`  // gas价格，最低为1
	Nonce     uint64 `json:"nonce"`     // 交易发起者nonce
	Sender    string `json:"sender"`    // 交易发起者公钥
	Body      struct {
		To    string `json:"to"`    // 交易接受者地址或合约地址，创建合约时为空
		Value string `json:"value"` // 交易金额
		Load  string `json:"load"`  // 合约负载，原生币OLO转账时为空
		Memo  string `json:"memo"`  // 备注 可为空
	} `json:"body"`
	Signature string `json:"signature"` // 交易签名
}
```

### 返回结果
```
交易hash
```

### 注意事项
- 交易应离线签名，并计算hash
- 调用接口广播后，应使用hash查询交易是否上链

### 示例
> POST http://127.0.0.1:10000/v2/contract/transactions
```
{
	"gasLimit": 10000000,
	"gasPrice": 1,
	"nonce": 2,
	"sender": "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8",
	"body": {
		"to": "0xe1066eBcFC8fbD7172886F15F538b63804676A74",
		"value": "0",
		"load": "8262963b0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000000066c6967616e670000000000000000000000000000000000000000000000000000"
	},
	"signature": "1af3046c03d8bbc423d60587533e6c019671f43f47000517cec9fa52253f9872588c6732bb3686dcd81b5ef368cf98fe614f8f008fde10ed2dd680bbbf01e46301"
}
```
应答
```
{
	"isSuccess": true,
	"message": "success",
	"result": {
		"code": 0,
		"gasUsed": 21000,
		"logs": "null",
		"ret": "",
		"tx": "0xe70b333dee039e3b11c13b591a59b87928932fe48ebdcf348a34aac4b1bf01bc"
	}
}
```

## 3.2 部署合约（私钥）
### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/contract/deploy](#)
> 可以使用3.3离线签名方法替代本方法
### 请求参数
```
type ContractDeployTx struct {
	GasLimit uint64 `json:"gasLimit"` // gas限额
	GasPrice string `json:"gasPrice"` // gas价格
	Sender   string `json:"sender"`   // 交易发起者公钥
	Privkey  string `json:"privkey"`  // 交易发起者私钥
	Value    string `json:"value"`    // 交易金额，通常为0
	Payload  string `json:"payload"`  // 合约部署字节码
	Memo     string `json:"memo"`     // 备注
}
```

### 返回结果
```
交易hash
```

### 示例
> POST http://127.0.0.1:10000/v2/contract/deploy
```
{
	"gasLimit": 10000000,
	"gasPrice": 1,
	"sender": "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8",
	"privkey": "7fffe4e426a6772ae8a1c0f2425a90fc6320d23e416fb6d83802889fa846faa2",
	"value": "0",
	"payload": "6060604052341561000f57600080fd5b6103a98061001e6000396000f30060606040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635a9b0b89146100515780638262963b146100e6575b600080fd5b341561005c57600080fd5b61006461014c565b6040518080602001838152602001828103825284818151815260200191508051906020019080838360005b838110156100aa57808201518184015260208101905061008f565b50505050905090810190601f1680156100d75780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34156100f157600080fd5b61014a600480803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919080359060200190919050506101fe565b005b6101546102c4565b600080600154818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156101ef5780601f106101c4576101008083540402835291602001916101ef565b820191906000526020600020905b8154815290600101906020018083116101d257829003601f168201915b50505050509150915091509091565b81600090805190602001906102149291906102d8565b50806001819055507f010becc10ca1475887c4ec429def1ccc2e9ea1713fe8b0d4e9a1d009042f6b8e82826040518080602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561028557808201518184015260208101905061026a565b50505050905090810190601f1680156102b25780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15050565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061031957805160ff1916838001178555610347565b82800160010185558215610347579182015b8281111561034657825182559160200191906001019061032b565b5b5090506103549190610358565b5090565b61037a91905b8082111561037657600081600090555060010161035e565b5090565b905600a165627a7a7230582066712744fea9374d65bc55f9ba7239759f8c88bd6f7c8439efa3f3555e1fad090029"
}
```
应答
```
{
	"isSuccess": true,
	"message": "success",
	"result": {
		"code": 0,
		"gasUsed": 1342466,
		"logs": "null",
		"ret": "606060405236156100d9576000357c......",
		"tx": "0xad54b931f6973905e0c9dd2d3c574bb42a638a002c611e6e93e3e94cdfddb888"
	}
}
```

## 3.3 执行合约（私钥）
### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/contract/invoke](#)
> > 可以使用3.3离线签名方法替代本方法

### 请求参数
```
type ContractInvokeTx struct {
   GasLimit        uint64 `json:"gasLimit"` //
   GasPrice        uint64 `json:"gasPrice"` //
   Sender          string `json:"sender"`   // pubkey
   Privkey         string `json:"privkey"`  // privkey
   Value           string `json:"value"`    // value
   ContractAddress string `json:"contract"` // contract address
   Payload         string `json:"payload"`  // abi.pack(function+参数)
   Memo            string `json:"memo"`     // 备注
}
```

### 返回结果
```
交易hash
```

### 示例
> POST http://127.0.0.1:10000/v2/contract/invoke
```
{
	"gasLimit": 10000000,
	"gasPrice": 1,
	"sender": "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8",
	"privkey": "7fffe4e426a6772ae8a1c0f2425a90fc6320d23e416fb6d83802889fa846faa2",
	"value": "0",
    "contract": "0xe1066eBcFC8fbD7172886F15F538b63804676A74",
	"payload": "8262963b0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000000066c6967616e670000000000000000000000000000000000000000000000000000"
}
```
应答
```
{
	"isSuccess": true,
	"message": "success",
	"result": {
		"code": 0,
		"gasUsed": 1342466,
		"logs": "null",
		"ret": "606060405236156100d9576000357c......",
		"tx": "0xad54b931f6973905e0c9dd2d3c574bb42a638a002c611e6e93e3e94cdfddb888"
	}
}
```

## 3.4 查询合约(私钥)
### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/contract/call](#)

### 请求参数
```
type ContractCallTx struct {
	GasLimit        uint64 `json:"gasLimit"` // gas限额
	GasPrice        string `json:"gasPrice"` // gas价格
	Sender          string `json:"sender"`   // 交易发起者公钥
	Privkey         string `json:"privkey"`  // 交易发起者私钥
	Value           string `json:"value"`    // 金额，通常为0
	ContractAddress string `json:"contract"` // 合约地址
	Payload         string `json:"payload"`  // 负载数据，abi.pack(function+参数) hex编码字符串
}
```

### 返回结果
```
type EvmCallResult = struct {
	Ret     string `json:"ret"`
	GasUsed uint64 `json:"gasUsed"`
}
```

### 示例
> POST http://127.0.0.1:10000/v2/contract/call
```
{
	"gasLimit": 10000000,
	"gasPrice": 1,
	"sender": "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8",
	"privkey": "7fffe4e426a6772ae8a1c0f2425a90fc6320d23e416fb6d83802889fa846faa2",
	"value": "0",
    "contract": "0xe1066eBcFC8fbD7172886F15F538b63804676A74",
	"payload": "5a9b0b89"
}
```
应答
```
{
	"isSuccess": true,
	"message": "",
	"result": {
		"code": 0,
		"gasUsed": 22707,
		"logs": null,
		"ret": "00000000000000000000000000000000000000000000000000470de28b962624",
		"tx": ""
	}
}
```

## 3.5 查询合约（签名）
### 请求说明
> 请求方式：POST<br>
> 请求URL ：[/v2/contract/query](#)

### 请求参数
```
type SignedEvmTx struct {
	CreatedAt uint64 `json:"createdAt"` // 时间戳，可选
	GasLimit  uint64 `json:"gasLimit"`  // gas限额
	GasPrice  string `json:"gasPrice"`  // gas价格，最低为1
	Nonce     uint64 `json:"nonce"`     // 交易发起者nonce
	Sender    string `json:"sender"`    // 交易发起者公钥
	Body      struct {
		To    string `json:"to"`    // 交易接受者地址或合约地址
		Value string `json:"value"` // 交易金额
		Load  string `json:"load"`  // 合约负载，普通原声币转账时为空
		Memo  string `json:"memo"`  // 备注
	} `json:"body"`
	Signature string `json:"signature"` // 交易签名
}
```

### 返回结果
```
type EvmCallResult = struct {
	Ret     string `json:"ret"`
	GasUsed uint64 `json:"gasUsed"`
}
```

### 示例
> POST http://127.0.0.1:10000/v2/contract/query
```
{
	"mode": 0,
	"createdAt": 1579070553927725700,
	"gasLimit": 1000000,
	"gasPrice": "1",
	"nonce": 58,
	"sender": "03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8",
	"body": {
		"to": "0xe1066eBcFC8fbD7172886F15F538b63804676A74",
		"value": "0",
		"load": "0x70a082310000000000000000000000000f508f143e77b39f8e20dd9d2c1e515f0f527d9f",
		"memo": "gosdk-v0.0.3"
	},
	"signature": "0xc9bbc3b8baad37fea40ca7582c7771e7c6c7acdcc3ddd9e06766408d63d08ece2a719815ace08641199763cbe197db48a06a4db80d0e121e4500d9898c8bc61500"
}
```
应答
```
{
	"isSuccess": true,
	"message": "",
	"result": {
		"code": 0,
		"gasUsed": 22707,
		"logs": null,
		"ret": "00000000000000000000000000000000000000000000000000470de28b962624",
		"tx": ""
	}
}
```


# 4. v3 新接口
## 4.1 查询所有账本
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/ledgers](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |N
address     |string          |用户地址或节点地址                   |N

### 返回结果
```
type V3Ledger struct {
	Id         uint64    `db:"id" json:"id"`               // 数据库自增id
	Height     int64     `db:"height" json:"height"`       // 区块高度
	BlockHash  string    `db:"blockHash" json:"blockHash"` // 区块hash
	BlockSize  int       `db:"blockSize" json:"blockSize"` // 区块大小：字节
	Validator  string    `db:"validator" json:"validator"` // 区块验证者节点地址
	TxCount    int64     `db:"txCount" json:"txCount"`     // 区块交易数
	GasLimit   int64     `db:"gasLimit" json:"gasLimit"`   // 区块gas限额之和
	GasUsed    int64     `db:"gasUsed" json:"gasUsed"`     // 区块所有交易消耗gas之和
	GasPrice   string    `db:"gasPrice" json:"gasPrice"`   // 区块交易平均gas价格，可能是小数
	CreatedAt  time.Time `db:"createdAt" json:"createdAt"` // 区块时间
}
```
### 示例


## 4.2 指定高度查询账本
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/ledgers/:height](#)

### 请求参数path
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
height       |int           |高度                    |Y

### 返回结果
同4.1

## 4.3 查询所有transaction
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/transactions](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |N
address     |string          |用户地址或节点地址                   |N

### 返回结果
```
type V3Transaction struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Hash      string    `db:"hash" json:"hash"`           // 交易hash
	Height    int64     `db:"height" json:"height"`       // 区块高度
	Typei     int       `db:"typei" json:"typei"`         // 交易类型
	Types     string    `db:"types" json:"types"`         // 交易类型
	Sender    string    `db:"sender" json:"sender"`       // 交易发起者地址
	Nonce     int64     `db:"nonce" json:"nonce"`         // 交易发起者nonce
	Receiver  string    `db:"receiver" json:"receiver"`   // 交易接受者地址
	Value     string    `db:"value" json:"value"`         // 交易金额
	GasLimit  int64     `db:"gasLimit" json:"gasLimit"`   // gas限额
	GasUsed   int64     `db:"gasUsed" json:"gasUsed"`     // gas使用量
	GasPrice  string    `db:"gasPrice" json:"gasPrice"`   // gas价格
	Memo      string    `db:"memo" json:"memo"`           // 备注
	Payload   string    `db:"payload" json:"payload"`     // 负载
	Events    string    `db:"events" json:"events"`       // 交易事件
	Codei     uint32    `db:"codei" json:"codei"`         // 失败代码
	Codes     string    `db:"codes" json:"codes"`         // 失败原因
	CreatedAt time.Time `db:"createdAt" json:"createdAt"` // 区块时间
}
```
Types       |Typei       |    备注
------------|-----------|-----------
TxTagAppInit |0|硬分叉账户初始化（忽略）
TxTagTinInit |256 |硬分叉tin初始化（忽略）
TxTagAppOLO  |257 |v1交易（废弃）
TxTagAppEvm  |513 |evm交易
TxTagAppFee  |769 |手续费交易（废弃）
TxTagAppBatch  |1025 |批量交易
TxTagNodeDelegate  |258 |节点抵押交易
TxTagUserDelegate  |514 |用户抵押交易
TxTagAppMgr  |65535 |链维护交易（忽略）

## 4.4 指定交易hash查询transaction
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/transactions/:txhash](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
txhash       |string           |交易hash                    |Y
### 返回结果
同4.3

## 4.5 指定高度查询transaction
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/ledgers/:height/transactions](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
height       |int64           |高度                    |Y
### 返回结果
同4.3

## 4.6 指定账户地址查询transaction
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/accounts/:address/transactions](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
address       | string          |账户地址                    |Y
### 返回结果
同4.3


## 4.7 查询所有payment
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/payments](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |Y
address     |string          |用户地址或节点地址                   |N

### 返回结果
```
type V3Payment struct {
	Id        uint64    `db:"id" json:"id"`               // 数据库自增id
	Hash      string    `db:"hash" json:"hash"`           // 交易hash
	Height    int64     `db:"height" json:"height"`       // 区块高度
	Idx       int       `db:"idx" json:"idx"`             // 交易索引
	Sender    string    `db:"sender" json:"sender"`       // 转账发起方地址
	Receiver  string    `db:"receiver" json:"receiver"`   // 转账接受方地址
	Symbol    string    `db:"symbol" json:"symbol"`       // 币种，原生币为“OLO”
	Contract  string    `db:"contract" json:"contract"`   // 合约地址，原生币为空或全零黑洞地址
	Value     string    `db:"value" json:"value"`         // 交易金额
	CreatedAt time.Time `db:"createdAt" json:"createdAt"` // 区块时间
}
```


## 4.8 指定交易hash查询payment
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/transactions/:txhash/payments](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |Y
address     |string          |用户地址或节点地址                   |N

### 返回结果
同4.7

## 4.9 指定高度查询payment
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/ledgers/:height/payments](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |Y
address     |string          |用户地址或节点地址                   |N

### 返回结果
同4.7

## 4.10 指定账户地址查询payment
### 请求说明
> 请求方式：GET<br>
> 请求URL ：[/v3/accounts/:address/payments](#)

### 请求参数query
字段       |字段类型       |字段说明                                |是否有效
------------|-----------|-----------|-------------------------
cursor       |int           |游标，从0开始                    |Y
limit          |int              |每页条数（最大200）                   |Y
order        |string        |排序方式(asc/desc)                   |Y
begin        |int              |开始时间戳                   |Y
end           |int              |结束时间戳                   |Y
height       |int64          |区块高度                   |N
symbol      |string          |币种，例如"OLO"，“ABC”                   |Y
address     |string          |用户地址或节点地址                   |N
### 返回结果
同4.7

