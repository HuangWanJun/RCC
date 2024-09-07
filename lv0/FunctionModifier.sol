// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

//Solidity 的修饰器是一种特殊的函数，它可以被用来修饰，
//改变函数的行为，修饰器可以在函数执行前进行预处理和验证操作。

contract FunctionModifier {
    // We will use these variables to demonstrate how to use
    // modifiers.
    address public owner;
    uint256 public x = 10;
    bool public locked;

     constructor() {
        // Set the transaction sender as the owner of the contract.
        owner = msg.sender;
    }

     // 检查用户必须一致，把检查抽取出来成为一个修饰器
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        // Underscore is a special character only used inside
        // a function modifier and it tells Solidity to
        // execute the rest of the code.
        _;
    }

 // 地址检查不能为0，把检查抽取出来成为一个修饰器
    modifier validAddress(address _addr) {
        require(_addr != address(0), "Not valid address");
        _;
    }
     // 使用多个修饰器，从左到右
    function changeOwner(address _newOwner) public onlyOwner validAddress(_newOwner)
    {
        owner = _newOwner;
    }

      modifier noReentrancy() {
        require(!locked, "No reentrancy"); //执行函数主体之前执行，异常处理会提前退出
        locked = true;  //执行函数主体之前执行
        _;              //代表函数主体
        locked = false; //执行函数主体之后执行
    }

    function decrement(uint256 i) public noReentrancy {
        x -= i;

        if (i > 1) {
            decrement(i - 1);
        }
    }
}
