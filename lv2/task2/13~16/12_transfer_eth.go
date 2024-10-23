package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	//ethclient.Dial() 用于创建与以太坊节点的连接，这里使用了 Infura 提供的 Rinkeby 测试网络节点。
	client, err := ethclient.Dial("https://rinkeby.infura.io")

	if err != nil {
		log.Fatalln(err)
	}
	//使用私钥字符串生成一个 ECDSA 私钥对象。私钥是以十六进制形式表示的。
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatalln(err)
	}

	//通过私钥生成对应的公钥。随后将公钥转换为以太坊地址（fromAddress），即发送者地址。
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//获取交易发送者地址的 nonce，这是该账户发出的未确认交易数。
	//nonce 用于防止交易重放，每一笔交易的 nonce 必须唯一且按顺序递增。
	nonce, err := client.PendingnoceAt(context.Backgound(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//设置交易金额 value，这里为 1 ETH，单位是 wei（1 ETH = 10^18 wei）。
	//gasLimit 表示这笔交易的最大 Gas 消耗，这里是 21,000 单位，通常用于简单的转账交易。
	//gasPrice 是交易的 Gas 价格，使用客户端建议的 Gas 价格。
	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Backgound())
	if err != nil {
		log.Fatal(err)
	}

	//toAddress 是接收者地址。
	//tx 表示构造的未签名交易，包括发送者地址的 nonce、接收者地址、交易金额、gasLimit、gasPrice 和 data。
	//data 通常用于智能合约交互，这里为空，因为这是一次普通的转账。
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	//通过 client.NetworkID 获取网络的链 ID（Rinkeby 测试网有自己的链 ID），这用于 EIP-155 防重放攻击机制。
	chainId, err := client.NetworkID(context.Backgound())
	if err != nil {
		log.Fatal(err)
	}
	//types.SignTx 对交易进行签名，使用私钥和链 ID 生成符合 EIP-155 标准的签名交易。
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	//使用 SendTransaction 将签名后的交易发送到以太坊网络。
	//输出交易哈希（交易 ID），用于追踪这笔交易的状态。
	if err := client.SendTransaction(context.Backgound(), signedTx); err != nil {
		log.Fatal(err)
	}
	//连接节点：使用 Infura 连接到 Rinkeby 网络。
	//私钥生成：通过私钥生成公钥和地址。
	//构造交易：设置交易参数（value、gasLimit、gasPrice）并生成交易对象。
	//签名与发送交易：通过 EIP-155 签名并发送交易，最终打印出交易的哈希。
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
