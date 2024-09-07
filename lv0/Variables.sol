// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

//Solidity 一共有三种数据位置，指定了变量数据位置。它们分别为：
// storage （数据会被存储在链上，是永久记录的）
// memory （数据存储在内存，是易失的） - 本地变量
// calldata （数据会被存在一个专门存放函数参数的地方）

contract Variables {
   //状态变量存储在区块链
    string public text = "Hello";
    uint256 public num = 123;

    function doSomething() public {
        // 本地变量不在区块链
        uint256 i = 456;

        //全局变量提供区块链的信息
        uint256 timestamp = block.timestamp; // Current block timestamp
        address sender = msg.sender; // address of the caller
    }
}