package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
)

// 使用原始交易（rawTx）字符串通过以太坊客户端将其发送到以太坊网络。
func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	rawTx := "f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"
	//这行代码将十六进制格式的交易字符串转换为字节数组，因为 RLP 编码后的交易数据通常以字节形式表示。
	rawTxBytes, err := hex.DecodeString(rawTx)

	//解码 RLP 并构造交易对象
	//tx := new(types.Transaction)：创建一个新的交易对象 tx，用来存储解码后的交易。
	tx := new(types.Transaction)
	//rlp.DecodeBytes(rawTxBytes, &tx)：使用 RLP 解码器将字节数组 rawTxBytes 解码成一个 types.Transaction 结构。
	//这将还原出原始的交易对象，包含所有交易细节。
	rlp.DecodeBytes(rawTxBytes, &tx)

	//SendTransaction 函数将已解码的 tx 发送到网络中。它会将交易广播到连接的以太坊节点，随后该节点会将交易发布到整个以太坊网络。
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0xc429e5f128387d224ba8bed6885e86525e14bfdc2eb24b5e9c3351a1176fd81f
}
