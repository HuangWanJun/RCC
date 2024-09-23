// 任何人都可以发送金额到合约
// 只有 owner 可以取款
// 3 种取钱方式

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.14;

contract EtherWallet{
    address payable public immutable owner;
    event Log(string funName,address from,uint256 value,bytes data);

    constructor() {
        owner = payable (msg.sender);
    }

    receive() external payable { 
        emit Log("receive", msg.sender, msg.value, "");
    }

    function withdrawal() external {
        require(msg.sender == owner, "not owner");
        payable (msg.sender).transfer(100);
    }

    function withdrawal2() external {
        require(msg.sender == owner,"Not owner");
        bool success = payable (msg.sender).send(200);
        require(success,"send failed");
    }

    function withdrawal3() external  {
        require(msg.sender == owner,"not owen");
        (bool success, ) = msg.sender.call{value: address(this).balance}("");
        require(success,"call failed");

    }

    function getBalance() external view returns(uint256){
        return address(this).balance;
    }
}