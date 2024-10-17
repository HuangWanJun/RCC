package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

// 包括生成私钥、公钥、计算钱包地址以及使用 Keccak256 哈希函数来验证钱包地址
func main() {

	//这里调用 crypto.GenerateKey() 来生成一对 ECDSA（椭圆曲线数字签名算法）的公私钥。
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	//crypto.FromECDSA(privateKey) 将私钥转换为字节数组。
	//hexutil.Encode(privateKeyBytes)[2:] 将字节数组转成十六进制表示，并去掉开头的 0x。
	//结果输出生成的私钥（以十六进制形式）
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	//privateKey.Public() 从私钥生成对应的公钥。
	//使用类型断言将 publicKey 转换为 *ecdsa.PublicKey 类型，如果类型不匹配程序会终止。
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	//crypto.FromECDSAPub(publicKeyECDSA) 将公钥转换为字节数组。
	//hexutil.Encode(publicKeyBytes)[4:] 将字节数组转为十六进制表示，去掉 0x 前缀和前两个字节（通常是 04，表示非压缩格式的公钥）。
	//结果输出公钥的十六进制形式。
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7d

	//crypto.PubkeyToAddress(*publicKeyECDSA) 通过公钥生成以太坊地址。这个过程使用公钥的哈希值的后 20 个字节作为地址。
	//address 是十六进制形式的以太坊地址，结果打印出地址。
	//第一种方式：直接通过公钥生成以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	//sha3.NewLegacyKeccak256() 创建一个 Keccak256 哈希函数的实例，这个哈希函数用于生成以太坊地址。
	//hash.Write(publicKeyBytes[1:]) 跳过公钥字节数组的第一个字节（即 04，表示非压缩格式），将剩余的字节写入哈希函数。
	//hexutil.Encode(hash.Sum(nil)[12:]) 计算哈希值并取最后 20 个字节（以太坊地址是公钥的 Keccak256 哈希值的最后 20 个字节）。
	//输出的哈希值应与上面生成的以太坊地址相匹配。

	//第二种方式：手动计算公钥哈希并提取地址
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e
}
