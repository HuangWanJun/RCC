// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// 一个基本类型合约的测试例子
contract Primitives {
    bool public boo = true;

    // uint 表示无符号整型
    uint8 public u8 = 1; // 8位
    uint256 public u256 = 456; //256位
    uint public u = 123; // uint默认256位

    // int 表示有符号整型
    int8 public i8 = -1;// 8位
    int256 public i256 = 456; //256位
    int256 public i = -123; // int默认256位

    // int类型的最小、最大值
    int256 public minInt = type(int256).min;
    int256 public maxInt = type(int256).max;

    address public addr = 0xCA35b7d915458EF540aDe6068dFe2F44E8fa733c;

    bytes1 a = 0xb5; //  [10110101]
    bytes1 b = 0x56; //  [01010110]
    bytes2 c = 0xb556;// [10110101 01010110]
    
    // 两种动态字节数组
    bytes d = new bytes(10); 
    string e = string(d); // 使用string()函数转换

    // 即使没有赋值也有默认值
    bool public defaultBoo; // false
    uint256 public defaultUint; // 0
    int256 public defaultInt; // 0
    address public defaultAddr; // 0x0000000000000000000000000000000000000000
}