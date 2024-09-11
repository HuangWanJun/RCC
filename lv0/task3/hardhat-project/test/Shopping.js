/*
 * @Author: Samuel Huang samuel.huang@tonghaifinancial.com
 * @Date: 2024-09-08 14:21:01
 * @LastEditors: Samuel Huang samuel.huang@tonghaifinancial.com
 * @LastEditTime: 2024-09-08 14:46:27
 * @FilePath: /hardhat-project/test/Shopping.js
 * @Description: 
 * 
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
 */


const {expect} = require("chai");
const hre = require("hardhat");
describe("Shopping", function () {
  let shoppingContract;
  before(async()=>{
     // ⽣成合约实例并且复⽤ 
     shoppingContract = await hre.ethers.deployContract("Shopping",[]);
  })

  it("should return the status Pending", async function () {
    
    expect(await shoppingContract.Status()).to.equal("Pending");
  });

  it("should return the status Shipped", async () => {
     await shoppingContract.Shopped();
     expect(await shoppingContract.Status()).to.equal("Shopped");
  });

  it("should return the status Delivered", async () => {
    await expect(shoppingContract.Devliverd()).to.emit(shoppingContract, "LogNewAlert")
    .withArgs("Your food order has arrived");
  });

});