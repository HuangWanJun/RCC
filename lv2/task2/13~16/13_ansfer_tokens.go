package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
)

// 转移 ERC-20 代币。
func main() {

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	privatekey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privatekey.PublicKey
	publickeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}
	fromAddress := crypto.PubkeyToAddress(*publickeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	//value 设为 0，因为这是 ERC20 代币转账，不涉及直接发送 ETH。
	value := big.NexInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HextToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	toKenAddress := common.HextToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	//构造 transfer 方法的调用数据
	//构造代币合约中 transfer 方法的函数签名 transfer(address,uint256)，
	//并通过 Keccak-256 哈希得到该方法的 methodID（前四个字节），
	//这是调用智能合约时用于标识方法的标识符。
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	//将接收者地址 toAddress 和转账数额（1000 个代币，转换为最小单位）进行左填充，确保每个参数占用 32 字节。
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddingAddress)) // 0x00000000000000000000000028b149020d2152179873ec60bed6bf7cd705775d

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	//将方法 ID、地址和金额拼接在一起，组成交易数据 data。这是传递给智能合约的完整数据。
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethclient.CallMsg{

		To: &toKenAddress,

		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit)

	//创建新的交易，nonce 是交易的顺序编号，tokenAddress 是代币合约地址，data 是调用合约的方法及参数。
	//通过链 ID 使用 EIP-155 签名机制签名交易，以确保交易只在特定网络上有效。
	tx := types.NewTransaction(nonce, toKenAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privatekey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // 0x9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0c8c0a9f0

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
