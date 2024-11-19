// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 导入 OpenZeppelin 提供的库和模块
/*
1. IERC20：ERC20 标准接口，用于与 ERC20 代币交互。
2. SafeERC20：安全的 ERC20 操作，避免直接调用 transfer 等方法导致的失败。
3. Address：地址工具类，提供安全的地址相关操作。
4. Math：数学运算工具。
5. Initializable：提供初始化逻辑，用于可升级合约。
6. UUPSUpgradeable：支持 UUPS 升级的合约基类。
7. AccessControlUpgradeable：基于角色的访问控制模块。
8. PausableUpgradeable：支持暂停功能的模块。
*/

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

// RCCStake 合约的功能：
// 1. 用户可以质押 ERC20 代币获取奖励。
// 2. RCC 奖励分配按质押比例计算。
// 3. 支持管理员设置池、暂停提现或领取奖励。
// 4. 提供了可升级性和权限控制。

contract RCCStake is
Initializable,  // 可升级合约的初始化模块
UUPSUpgradeable, // 支持 UUPS 升级
PausableUpgradeable, // 支持暂停功能
AccessControlUpgradeable // 基于角色的权限管理
{
    using SafeERC20 for IERC20; // 安全的 ERC20 操作
    using Address for address; // 提供地址相关的工具函数
    using Math for uint256;    // 提供数学运算工具

    // ************************************** 常量定义 **************************************

    /*
    定义的角色权限：
    - ADMIN_ROLE：用于执行管理功能。
    - UPGRADE_ROLE：用于合约升级操作。
    */
    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");

    // 定义 ETH 专属池 ID，所有与 ETH 池相关的操作将使用此 ID
    uint256 public constant ETH_PID = 0;
    // ************************************** 数据结构 **************************************

    /*
   Pool 结构体：
   - 表示每个质押池的配置与状态。
   */

    struct Pool {
        address stTokenAddress;       // 质押代币地址
        uint256 poolWeight;           // 池权重，用于奖励分配比例
        uint256 lastRewardBlock;      // RCC 最后一次奖励分配的区块号
        uint256 accRCCPerST;          // 每单位质押代币的累计 RCC 奖励（放大）
        uint256 stTokenAmount;        // 当前池中质押的代币总量
        uint256 minDepositAmount;     // 用户最低质押金额
        uint256 unstakeLockedBlocks;
    }
    // 解押请求锁定的区块数

    /*
UnstakeRequest 结构体：
- 记录用户的解押请求，用于延迟解押。
*/
    struct UnstakeRequest {
        uint256 amount;       // 解押的代币数量
        uint256 unlockBlocks; // 解锁该请求的区块号
    }

    /*
   User 结构体：
   - 用户的质押状态和奖励状态。
   */
    struct User {
        uint256 stAmount;            // 用户已质押的代币总量
        uint256 finishedRCC;         // 已完成的 RCC 奖励
        uint256 pendingRCC;          // 待领取的 RCC 奖励
        UnstakeRequest[] requests;   // 用户的解押请求数组,解质押请求列表，每个请求包含解质押数量和解锁区块。
    }
    // ************************************** 状态变量 **************************************

    /*
   RCC 奖励分配规则：
   - startBlock：奖励开始分配的区块号。
   - endBlock：奖励结束分配的区块号。
   - rccPerBlock：每个区块分配的 RCC 奖励数量。
   */
    uint256 public startBlock;
    uint256 public endBlock;
    uint256 public rccPerBlock;
    /*
   功能控制开关：
   - withdrawPaused：是否暂停提现。
   - claimPaused：是否暂停领取奖励。
   */
    bool public withdrawPaused;
    bool public claimPaused;
    // RCC 代币的地址
    IERC20 public RCC;

    // 总的池权重，用于奖励分配的权重基准
    uint256 public totalPoolWeight;

    // 质押池数组，每个池保存为 Pool 结构体
    Pool[] public pool;
    // 用户信息映射，按照 [质押池ID => 用户地址 => 用户状态] 的结构存储
    mapping(uint256 => mapping(address => User)) public user;


    // ************************************** 事件 **************************************
 event SetRCC(IERC20 indexed RCC);
    event PauseWithdraw();
    event UnpauseWithdraw();
    event PauseClaim();
    event UnpauseClaim();
    event SetStartBlock(uint256 indexed startBlock);
    event SetEndBlock(uint256 indexed endBlock);
    /*
    管理操作相关事件：
    - 用于记录重要操作。
    */

    // ************************************** 修饰符 **************************************

    /*
       checkPid 修饰符：
       - 验证给定的质押池 ID 是否有效。
       */

    modifier checkPid(uint256 _pid) {
        require(_pid < pool.length, "Invalid pid");
        _;
    }

/*
whenNotClaimPaused 修饰符：
- 确保领取奖励操作未暂停。
*/
    modifier whenNotClaimPaused() {
        require(!claimPaused, "Claim paused");
        _;
    }

/*
   whenNotWithdrawPaused 修饰符：
   - 确保提现操作未暂停。
   */
    modifier whenNotWithdrawPaused(){
        require(!claimPaused, "withdraw is paused");
        _;
    }

// ************************************** 初始化函数 **************************************
/*
    initialize：合约初始化函数。
    - 参数：
      - _RCC：RCC 代币的地址。
      - _startBlock：奖励分配的开始区块。
      - _endBlock：奖励分配的结束区块。
      - _rccPerBlock：每个区块分配的 RCC 奖励数量。
    - 初始化角色权限和奖励机制。
    */
    function initialize(
        address _RCC,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _rccPerBlock
    ) public initializer
    {
        require(_startBlock <= _endBlock && _rccPerBlock > 0, "invalid parameters");

// 初始化访问控制与升级模块
        __AccessControl_init();
        __UUPSUpgradeable_init();

// 设置默认管理员角色
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);

// 设置 RCC 代币地址
        setRCC(_RCC);

// 初始化奖励分配参数
        startBlock = _startBlock;
        endBlock = _endBlock;
        rccPerBlock = _rccPerBlock;
    }

     /*
    _authorizeUpgrade：用于验证升级操作，仅允许具有升级权限的地址执行。
    */
     function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    {}
// ************************************** 管理功能 **************************************
    /*
    setRCC：设置 RCC 代币地址。
    - 仅限具有 ADMIN_ROLE 的地址调用。
    */
    function setRCC(IERC20 _RCC) public onlyRole(ADMIN_ROLE){
        RCC=_RCC;
        emit SetRCC(RCC);
    }

    /*
    pauseWithdraw / unpauseWithdraw：
    - 分别用于暂停/恢复提现操作。
    - 仅限具有 ADMIN_ROLE 的地址调用。
    */
    function pauseWithdraw() public onlyRole(ADMIN_ROLE){
        require(!withdrawPaused, "Withdraw is already paused");
        withdrawPaused = true;
        emit PauseWithdraw();
    }

    function unpauseWithdraw() public onlyRole(ADMIN_ROLE){
        require(withdrawPaused, "Withdraw is not paused");
        withdrawPaused = false;
        emit UnpauseWithdraw();
    }
    //暂停和取消暂停 Claim
    function pauseClaim() public onlyRole(ADMIN_ROLE)  {
        require(!claimPaused, "Claim is already paused");
        claimPaused = true;
        emit PauseClaim();
    }

    function unpauseClaim() public onlyRole(ADMIN_ROLE)  {
        require(claimPaused, "Claim is not paused");
        claimPaused =false;
        emit UnpauseClaim();
    }

    //设置区块范围
    function setStartBlock(uint256 _startBlock) public onlyRole(ADMIN_ROLE)  {
        require(_startBlock <= endBlock, "start block must be smaller than end block"));
        startBlock = _startBlock;
        emit SetStartBlock(_startBlock);
    }

    function setEndBlock(uint256 _endBlock) public onlyRole(ADMIN_ROLE)  {
        require(startBlock <= _endBlock,"start block must be smaller than end block");
        endBlock = _endBlock;
        emit SetEndBlock(_endBlock);
    }




}