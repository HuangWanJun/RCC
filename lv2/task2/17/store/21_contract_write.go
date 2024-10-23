package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	publickey := privateKey.Public()
	publickeyECDSA, ok := publickey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publickey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publickeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//创建交易授权对象
	//使用 bind.NewKeyedTransactor 创建一个授权对象 auth，该对象包含了交易的发送者身份信息（私钥）。
	//还设置了交易的 Nonce，Value（发送的以太币数量，这里为 0），GasLimit（最大 gas 量），以及 GasPrice。
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	//加载智能合约实例  地址
	address := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	//设置数据（调用智能合约中的 SetItem 方法）
	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))
	//准备好 key 和 value，并调用合约的 SetItem 方法，将这些数据写入合约存储中。
	//SetItem 方法会发送一笔交易，交易会被广播到网络中，tx.Hash().Hex() 打印出该交易的哈希值，方便追踪。
	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870

	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result[:])) // "bar"
}
