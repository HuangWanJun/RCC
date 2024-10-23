// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

/*
go get -u github.com/ethereum/go-ethereum
cd $GOPATH/src/github.com/ethereum/go-ethereum/
make
make devtools
*/

/*

//这个命令用于从 Solidity 源代码文件 Store.sol 中生成 ABI（Application Binary Interface）。
ABI 是一种描述智能合约接口的 JSON 数据结构，包含合约中的所有函数及其参数类型。
生成的 ABI 文件通常用于与智能合约进行交互。

solcjs --abi Store.sol


//这个命令用于编译 Store.sol 文件并生成合约的字节码（binary）。字节码是智能合约在以太坊虚拟机（EVM）上执行的实际代码。
//生成的字节码通常用于部署智能合约到以太坊网络。
solcjs --bin Store.sol

//这个命令使用 abigen 工具将前面生成的 ABI 和字节码文件转换为 Go 语言中的合约调用代码。
生成的 Go 文件（Store.go）将包含可以直接在 Go 代码中使用的合约函数的结构体和方法。
//这使得在 Go 项目中与 Solidity 智能合约进行交互变得更加简单和高效。
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=Store.go

//solcjs：用于编译 Solidity 合约文件，生成 ABI 和字节码。
//abigen：用于将 ABI 和字节码转换为 Go 语言代码，方便在 Go 项目中与智能合约进行交互。
*/

contract Store {
    event ItemSet(bytes32 key, bytes32 value);
    string public version;
    mapping (bytes32 => bytes32)public items;
    constructor(string memory _version){
        version = _version;
    }
    function setItem(bytes32 key, bytes32 value) public {
        items[key] = value;
        emit ItemSet(key, value);
    }
}
