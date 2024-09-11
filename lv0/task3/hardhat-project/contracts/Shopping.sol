pragma solidity ^0.8.7;

contract Shopping {
    enum ShoppingStatus {
        Pending,
        Shopped,
        Delivered
    }

    ShoppingStatus private status;
    event LogNewAlert(string descrption);

    constructor(){
        status = ShoppingStatus.Pending;
    }

    function Shopped()public{
        status = ShoppingStatus.Shopped;
        emit LogNewAlert("Your food order has been start to delivery");
    }

    function Devliverd() public {
        status = ShoppingStatus.Delivered;
        emit LogNewAlert("Your food order has arrived");
    }

    function getStatus(ShoppingStatus _status) 
        internal pure returns (string memory) {
        if(ShoppingStatus.Pending == _status) return "Pending";
        if(ShoppingStatus.Shopped == _status) return "Shopped";
        if(ShoppingStatus.Delivered == _status) return "Delivered";
        else return "";
    }

    function Status() public view returns (string memory) {
        ShoppingStatus _status = status;
        return getStatus(_status);
    }
}