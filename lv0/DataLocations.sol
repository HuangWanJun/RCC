// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// Solidity 要求在声明「引用类型」的时候必须要加上数据位置(data location)，有三种：
// 1、storage （数据会被存储在链上，是永久记录的，其生命周期与合约生命周期一致）
// 2、memory （数据存储在内存，是易失的，其生命周期与函数调用生命周期一致，函数调用结束数据就消失了）
// 3、calldata （与memory类似，数据会被存在一个专门存放函数参数的地方，与memory不同的是calldata数据是不可更改的。另外相比于memory，它消耗更少的Gas）

contract DataLocations {
    uint256[] public arr;
    mapping(uint256 => address) map;

    struct MyStruct{
        uint256 foo;
    }

     mapping(uint256 => MyStruct) myStructs;

     function f() public{
        _f(arr,map,myStructs[1]);
          // get a struct from a mapping
        MyStruct storage myStruct = myStructs[1];
            // create a struct in memory
        MyStruct memory myMemStruct = MyStruct(0);
     }
     
     function _f(
        uint256[] storage _arr,
        mapping(uint256 => address) storage _map,
        MyStruct storage _myStruct
     ) internal {
         // do something with storage variables
     }

      // You can return memory variables
    function g(uint256[] memory _arr) public returns (uint256[] memory) {
        // do something with memory array
    }

    function h(uint256[] calldata _arr) external {
        // do something with calldata array
    }
}