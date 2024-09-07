// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract Function{
     // 返回多个值
     function returnMany() public pure returns(uint256,bool,uint256){
        return (1,true,2);
     }

      function named() public pure returns (uint256 x, bool b, uint256 y) {
        return (1, true, 2);
    }

      // 返回值也可以被赋值
    function assigned() public pure returns (uint256 x, bool b, uint256 y) {
        x = 1;
        b = true;
        y = 2;
    }

     // 在呼叫另一個
    // 返回多個值的函數時，請使用析构赋值。
      function destructuringAssignments() public pure returns (uint256, bool, uint256, uint256, uint256)
      {
         (uint256 i, bool b, uint256 j) = returnMany();
         (uint256 x,, uint256 y) = (4, 5, 6);
          return (i, b, j, x, y);
      }

        function arrayInput(uint256[] memory _arr) public {}

    // Can use array for output
    // 数组也可以作为返回值
    uint256[] public arr;

    function arrayOutput() public view returns (uint256[] memory) {
        return arr;
    }
}