// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract SimpleStorage {
    // State variable to store a number
    uint256 public num;

    // You need to send a transaction to write to a state variable.
    function set(uint256 _num) public {
        num = _num;
    }
//view（视图）函数只能读取合约状态，不能修改合约状态。
//pure 纯函数不读不写，没有副作用。使用纯函数可以提高代码安全性，避免出现与预期不符的副作用。
   // 你可以免费读取状态变量
    function get() public view returns (uint256) {
        return num;
    }
}