package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)
	//在以太坊网络中，每个交易都有一个递增的 nonce 值。这个值确保交易的唯一性，并避免双重支出。这行代码从网络中获取了发送地址的当前 nonce 值。
	//in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	//将签名的交易转换为 RLP 编码
	ts := types.Transactions{signedTx}
	//ts.GetRlp(0)：将签名的交易转换为 RLP（Recursive Length Prefix）编码
	//，RLP 是以太坊中用于编码对象的标准格式。
	rawTxBytes := ts.GetRlp(0)
	//将 RLP 编码后的字节转换为十六进制字符串。
	rawTxHex := hex.EncodeToString(rawTxBytes)

	//这个字符串可以直接用于广播到以太坊网络。
	fmt.Printf(rawTxHex) // f86...772
}
