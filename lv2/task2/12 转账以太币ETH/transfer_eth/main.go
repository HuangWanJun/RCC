package main

import (
	"log"
)

// TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.
// 简单的以太坊交易
func main() {

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	//使用私钥字符串创建一个椭圆曲线数字签名算法 (ECDSA) 私钥对象。这个私钥是用于对交易进行签名的。
	//如果私钥格式不正确或加载失败，程序会终止。

	privateKey, err := crytop.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")

}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
