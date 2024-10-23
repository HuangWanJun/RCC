package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	store "store/contracts"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/eba825833e774dd582d26c47a78184bb")
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
	//  // 获取待处理的 nonce 值（即当前交易数量）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 创建一个带有私钥的交易授权者
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce)) // 设置 nonce
	auth.Value = big.NewInt(0)            // 设置发送的以太值（这里为0）
	auth.GasLimit = uint64(300000)        // 设置 gas 限制
	auth.GasPrice = gasPrice              // 设置 gas 价格
	//余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %s wei\n", balance.String())

	// 计算交易费用
	estimatedGasCost := new(big.Int).Mul(big.NewInt(int64(auth.GasLimit)), gasPrice)
	// 打印账户余额和估算的交易费用
	fmt.Printf("Balance: %s wei\n", balance.String())
	fmt.Printf("Estimated Gas Cost: %s wei\n", estimatedGasCost.String())

	if balance.Cmp(estimatedGasCost) < 0 {
		log.Fatal("Insufficient funds for gas * price")
	}

	// 部署合约，input 是构造函数的参数
	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err) // 如果部署合约失败，则记录错误并终止程序
	}
	// 输出合约地址和交易哈希
	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	// 保存合约实例，供后续使用
	_ = instance
}
