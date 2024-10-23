package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 订阅新区块
// 订阅新的区块头并处理这些区块头相关的信息。它会监听每个新生成的区块并输出区块的详细信息。以下是对代码的详细解析：
func main() {
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	//订阅新生成的区块头
	//make(chan *types.Header) 创建了一个接收 types.Header 的通道 headers，用于存储新的区块头。
	//SubscribeNewHead() 方法订阅新生成的区块头。当新的区块头被矿工生成并添加到链上时，
	//会触发这个订阅并将区块头发送到 headers 通道。
	//sub 是订阅对象，用于检查错误和取消订阅。
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	//进入监听循环
	for {
		select {
		case err := <-sub.Err():
			log.Println(err) //sub.Err()：如果订阅过程中发生错误，将终止程序并打印错误。
			//headers：当新的区块头到达时，从 headers 通道接收区块头并处理。
		case header := <-headers:
			// 处理新接收的区块头
			//header.Hash()：每个区块头都有一个唯一的哈希值，用于标识该区块。该行打印区块头的哈希值。
			fmt.Println(header.Hash().Hex())
			//client.BlockByHash()：通过区块头的哈希值获取完整的区块信息。区块头仅包含部分信息，完整区块还包括交易等数据。
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			// // 打印区块的哈希值
			fmt.Println(block.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			//// 打印区块的编号
			fmt.Println(block.Number().Uint64()) // 3477413
			// // 打印区块的时间戳
			fmt.Println(block.Time().Uint64()) // 1529525947
			//  // 打印区块的随机数（用于工作量证明）
			fmt.Println(block.Nonce()) // 130524141876765836
			//// 打印区块中的交易数量
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}
