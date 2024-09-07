// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// 事件 类似于log, 用来记录函数执行事件

contract Event {
    // 事件聲明
    // 最多可索引 3 個參數。
    // 索引參數可協助您依索引參數過濾日誌
    event Log(address indexed sender, string message);
    event AnotherLog();

    function test() public {
        emit Log(msg.sender, "Hello World!");
        emit Log(msg.sender, "Hello EVM!");
        emit AnotherLog();
    }
}
