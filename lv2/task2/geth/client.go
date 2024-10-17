package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"    //提供与以太坊地址和数据类型处理相关的功能。
	"github.com/ethereum/go-ethereum/ethclient" //Go-ethereum 的客户端库，用于与以太坊节点通信。
	"log"
	"math"
	"math/big"
)

func main() {

	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	//获取账户余额

	//将十六进制字符串转换为以太坊地址格式。
	account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	//client.BalanceAt：获取指定账户的余额。nil 表示获取最新区块的余额。
	//返回余额，单位是 wei（以太坊的最小单位，1 ETH = 10^18 wei）。
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	//获取特定区块的余额
	//创建一个 big.Int 对象，代表区块编号 5532993。
	blockNumber := big.NewInt(5532993)
	//client.BalanceAt：获取指定区块中账户的余额。
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balanceAt)

	//将余额转换为以太币单位
	//big.Float：用于表示浮点数，适合处理精度较大的数值。
	fbalance := new(big.Float)
	//balanceAt.String()：将 balanceAt 转换为字符串。
	fbalance.SetString(balanceAt.String())
	//math.Pow10(18)：计算 10^18，用于将 wei 转换为以太币单位。
	//ethValue：将 balanceAt 转换为以太币（ETH）。
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue)

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance) // 25729324269165216042

}

func getAddress() {

	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Println(address.Hex())
	//fmt.Println(address.Hash().Hex())
	fmt.Println(address.Bytes())
}
