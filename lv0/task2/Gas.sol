// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// 一ether等于 10(18) wei。
contract Gas {
    uint256 public i = 0;

    // gas用完交易会失败
    // 状态改变还原
    // 已花费的Gas不会退
    function forever() public {
        // Here we run a loop until all of the gas are spent
        // and the transaction fails
        while (true) {
            i += 1;
        }
    }
}