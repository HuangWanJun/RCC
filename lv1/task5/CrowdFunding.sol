//众筹合约

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract CrowdFunding {
    address public immutable beneficiary; //受益人
    uint256 public immutable fundingGoal; //筹资目标数量
    uint256 public fundingAmount; //当前的 金额

    mapping(address => uint256) public funders;
    // 可迭代的映射
    mapping(address => bool) private fundersInserted;
    address[] public funderskey; //length
    // 不用自销毁方法，使用变量来控制
    bool public AVALIABLED = true;

    // 部署的时候，写入受益人+筹资目标数量
    constructor(address beneficiary_, uint256 golal_) {
        beneficiary = beneficiary_;
        fundingGoal = golal_;
    }
    // 资助
    //      可用的时候才可以捐
    //      合约关闭之后，就不能在操作了
    function contribute() external payable {
        require(AVALIABLED,"CrowdFunding is closed");

        // 检查捐赠金额是否会超过目标金额
        uint256 potentialFundingAmout = fundingAmount + msg.value;
        uint256 refundAmount = 0;

        if(potentialFundingAmout > fundingGoal){
            refundAmount = potentialFundingAmout - fundingGoal;
            funders[msg.sender] += (msg.value - refundAmount);
            fundingAmount += (msg.value - refundAmount);
        }else{
            funders[msg.sender] += msg.value;
            fundingAmount += msg.value;
        }

         // 更新捐赠者信息
         if(!fundersInserted[msg.sender]){
            fundersInserted[msg.sender] = true;
            funderskey.push(msg.sender);
         }

         // 退还多余的金额
         if(refundAmount >0){
            payable (msg.sender).transfer(refundAmount);
         }
    }

     // 关闭
     function close() external returns (bool){
        //1.checking
        if(fundingAmount < fundingGoal){
            return false;
        }
        uint256 amount = fundingAmount;
        //2.edit
        fundingAmount = 0;
        AVALIABLED = false;
        //3.action
        payable (beneficiary).transfer(amount);
        return true;
     }

     function funderLenght() public  view returns (uint256){
        return funderskey.length;
     }


}
