// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// 构造函数是在创建合约时执行的可选函数。
// 构造函数是 Solidity 中的一种特殊函数。它的主要作用是对合约状态进行初始化。
// 构造函数为我们提供了把变量初始值参数化的途径。

contract X {
    string public name;
    constructor(string memory _name){
        name = _name;
    }

   
}

 contract Y {
    string public text;
    constructor(string memory _text) {
        text = _text;
    }
 }

 // 下面2种方法初始父合约
/// 1、在继承列表
contract B is X("input to x"), Y("input to Y"){}

contract C is X,Y{
    // 2、在构造函数
    constructor(string memory _name,string memory _text) X(_name) Y(_text) {
        
    }


}

// 父构造函数始终按继承顺序调用，而不是按照子合约构造函数列出的顺序。
// 构造函数调用的顺序
// Order of constructors called:
// 1. X
// 2. Y
// 3. D
contract D is X, Y {
    constructor() X("X was called") Y("Y was called") {}
}

// 构造函数调用的顺序
// Order of constructors called:
// 1. X
// 2. Y
// 3. E
contract E is X, Y {
    constructor() Y("Y was called") X("X was called") {}
}


