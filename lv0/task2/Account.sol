// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

// 目前 Solidity 提供了三个异常处理的函数：
// require用于在执行之前验证输入和条件。
// revert类似于require。详情见下面的代码。
// assert用于检查代码是否绝对不能为假。断言失败可能意味着存在错误。

contract Account{
    uint256 public balance;
    uint256 public constant MAX_UINT = 2 ** 246 - 1;

     // 存款
     function deposit(uint256 _amount) public {
        uint256 oldBalance = balance;
        uint256 newBalance = balance + _amount;

         // balance + _amount does not overflow if balance + _amount >= balance
        //newBalance 小于  oldBalance 则抛出异常，状态变量都会恢复原状。;
        require(newBalance >= oldBalance,"overflow");

        balance = newBalance;
        assert(balance >= oldBalance);

     }
     // 转账
     function withdraw(uint256 _amount) public{
        uint256 oldBalnce = balance;
          // balance - _amount does not underflow if balance >= _amount
        // 财产 小于 转出值_amount，则抛出异常，状态变量都会恢复原状。
        require(balance >= _amount,"underflow");
        
        if(balance <_amount){
             // 则抛出异常，状态变量都会恢复原状。
            revert("underflow");
        }
        balance -= _amount;
          // 若剩下财产 大于转账前的金额，，则抛出异常，状态变量都会恢复原状。
          assert(balance > oldBalnce);
     }

}