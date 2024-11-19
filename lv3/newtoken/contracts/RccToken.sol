pragma solidity ^0.8.27;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract RccToken is ERC20{
    constructor() ERC20("RCCToken", "RCC") {
        _mint(msg.sender, 2000000*1_000_000_000_000_000_000);
    } 
}