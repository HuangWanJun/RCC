package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}
	//使用 client.BlockByNumber 从以太坊区块链中获取区块号为 5671744 的区块。
	//如果获取区块失败，程序会输出错误并终止。
	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	//block.Transactions() 返回区块中的所有交易，通过遍历每个交易 tx，代码输出交易的关键信息，
	//如哈希值、金额、gas、gas 价格、nonce 值、附加数据和接收方地址。
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		fmt.Println(tx.Value().String())    // 10000000000000000
		fmt.Println(tx.Gas())               // 105000
		fmt.Println(tx.GasPrice().Uint64()) // 102000000000
		fmt.Println(tx.Nonce())             // 110644
		fmt.Println(tx.Data())              // []
		fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e

		//获取交易的发送者地址

		//它用于获取以太坊网络的 链 ID（chainID）。链 ID 是一个用来标识不同以太坊网络的数字，不同的网络有不同的链 ID
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		//通过获取链的 chainID 和使用 types.Sender 函数，代码可以推导出交易的发送者地址。
		//EIP-155 是以太坊的一项改进提案，用于解决跨链重放攻击的问题。
		//EIP-155 引入了链 ID 签名机制，将链 ID 添加到交易的签名数据中，以确保交易只能在特定链上有效。
		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println("sender", sender.Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		}

		//client.TransactionReceipt 根据交易哈希获取交易收据，其中包含交易的状态。
		//打印 receipt.Status，如果值为 1，则交易成功。
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status) // 1
	}

	//方法二，	//使用区块哈希 blockHash 获取该区块中的交易数量。
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	//通过交易索引，使用 client.TransactionInBlock 获取区块中的每笔交易并打印其哈希值。
	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	}
	//根据交易哈希获取交易
	///根据交易哈希 txHash，使用 client.TransactionByHash 获取特定交易，
	//并判断交易是否还在待处理状态 (pending)。
	txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(isPending)       // false

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
