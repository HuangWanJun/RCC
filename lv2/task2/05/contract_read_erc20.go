package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	token "readerc/contracts_erc20"
)

func main() {
	//ethclient.Dial 用于连接到以太坊节点。在这里，使用了 Cloudflare 的以太坊节点服务。
	//如果连接失败，程序会记录错误并终止。
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	// Golem (GNT) Address
	//common.HexToAddress 将十六进制字符串转换为以太坊地址，指定 Golem (GNT) 的合约地址。
	tokenAddress := common.HexToAddress("0xa74476443119A942dE498590Fe1f2454d7D4aC0d")
	//token.NewToken 创建了一个新的代币合约实例，以便调用代币的智能合约方法。如果创建失败，程序会记录错误并终止。
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	//查询用户的代币余额。
	//定义一个用户地址 address，使用 BalanceOf 方法查询该地址持有的代币余额。
	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	//&bind.CallOpts{} 用于提供调用选项（在这里为空，表示默认选项）。
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
	//将余额转换为人类可读的格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	//将余额转换为可读的代币格式。
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
}
