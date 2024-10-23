package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	store "store/contracts"
)

func main() {

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/eba825833e774dd582d26c47a78184bb")
	if err != nil {
		panic(err)
	}
	//合约地址
	address := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}
	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version) // "1.0"
}
