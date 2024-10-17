package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	createKs()
	//importKs()
}

func createKs() {

	//初始化一个新的密钥存储（keystore），存放在当前目录（"./"）。StandardScryptN 和 StandardScryptP 是标准的 Scrypt 加密参数，定义了密钥派生的安全性。
	ks := keystore.NewKeyStore("./", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret1"
	//使用指定的密码 "secret" 创建一个新的以太坊账户，keystore 会将私钥安全地保存，并生成一个新的以太坊地址。
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	//打印新创建的以太坊账户的地址，格式为十六进制。
	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
}

// 从已有的 keystore 文件中导入以太坊账户，并且可以选择性地删除文件。
func importKs() {

	//指定要导入的 keystore 文件的路径，文件名包含时间戳和以太坊地址。
	file := "./UTC--2024-10-17T09-29-59.261697000Z--01d5b3d4efa8491dd8cae1893c158da63876790a"
	ks := keystore.NewKeyStore("./", keystore.StandardScryptN, keystore.StandardScryptP)
	//读取指定的 keystore 文件内容，并存储在 jsonBytes 变量中。
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	password := "secret1"
	//导入 keystore 文件中的私钥，使用第一个 password 来解锁私钥，第二个 password 用于加密并保存私钥。
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {

		log.Fatal(err)
	}
	//打印导入的以太坊账户的地址。
	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
	//删除导入的 keystore 文件。
	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}
