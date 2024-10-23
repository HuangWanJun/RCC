package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// /查询区块
func main() {

	//ethclient.Dial 用于连接到指定的以太坊节点，这里使用了 Cloudflare 提供的公共以太坊节点 https://cloudflare-eth.com。
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	//client.HeaderByNumber 获取某个区块号对应的区块头信息。
	//nil 表示获取最新的区块头。
	//context.Background() 是传递给 API 的上下文信息，这里使用了默认的背景上下文。
	//如果获取区块头失败，程序会打印错误并退出。
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String()) // 5671744
	//blockNumber := big.NewInt(5671744) 创建一个表示区块号 5671744 的 big.Int 对象。
	//client.BlockByNumber 获取区块号为 5671744 的完整区块信息。
	//如果获取区块信息失败，程序会打印错误并退出。
	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	//block.Number().Uint64() 输出区块号（5671744）。
	//block.Time().Uint64() 输出区块时间戳（秒为单位），这里的时间戳为 1527211625，可以转换为具体日期时间。
	//block.Difficulty().Uint64() 输出区块的挖矿难度值（3217000136609065）。
	//block.Hash().Hex() 输出区块的哈希值（0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9）。
	//len(block.Transactions()) 输出区块中包含的交易数量（144 笔交易）。
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1527211625
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144

	//client.TransactionCount 获取指定区块的交易数量。
	//通过区块的哈希值来查询区块内的交易数量。
	//输出查询到的交易数量，这里为 144。
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count) // 144

}
