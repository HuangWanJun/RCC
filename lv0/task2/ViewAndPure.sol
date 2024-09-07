// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// Solidity 提供了三个状态可变性修饰符：
// view 函数只能查询合约状态，不能更改合约状态。简单来讲就是只读不写的
// pure 既不能查询，也不能修改函数状态。只能使用函数参数进行简单计算并返回结果
// payable 允许函数接受 Ether 转账。函数默认情况下是不能接受转账的，如果你需要接受转账，那么必须指定其为 payable


contract ViewAndPure {
    uint256 public x = 1;

    // Promise not to modify the state.
    /// view 函数只能查询合约状态，不能更改合约状态。简单来讲就是只读不写的
    function addToX(uint256 y) public view returns (uint256) {
        return x + y;
    }

    // Promise not to modify or read from the state.
    // pure 既不能查询，也不能修改函数状态。只能使用函数参数进行简单计算并返回结果
    function add(uint256 i, uint256 j) public pure returns (uint256) {
        return i + j;
    }
}