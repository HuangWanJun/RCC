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
    event setRccPerBlock(uint256 indexed rccPerBlock);
    event AddPool(address indexed stTokenAddress, uint256 indexed poolWeight, uint256 indexed lastRewardBlock, uint256 minDepositAmount, uint256 unstakeLockedBlocks);

    event UpdatePoolInfo(uint256 indexed poolId, uint256 indexed minDepositAmount, uint256 indexed unstakeLockedBlocks);

    event SetPoolWeight(uint256 indexed poolId, uint256 indexed poolWeight, uint256 totalPoolWeight);

    event UpdatePool(uint256 indexed poolId, uint256 indexed lastRewardBlock, uint256 totalRCC);

    event Deposit(address indexed user, uint256 indexed poolId, uint256 amount);

    event RequestUnstake(address indexed user, uint256 indexed poolId, uint256 amount);

    event Withdraw(address indexed user, uint256 indexed poolId, uint256 amount, uint256 indexed blockNumber);

    event Claim(address indexed user, uint256 indexed poolId, uint256 rccReward);

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
    function setRCC(IERC20 _RCC) public onlyRole(ADMIN_ROLE) {
        RCC = _RCC;
        emit SetRCC(RCC);
    }

    /*
    pauseWithdraw / unpauseWithdraw：
    - 分别用于暂停/恢复提现操作。
    - 仅限具有 ADMIN_ROLE 的地址调用。
    */
    function pauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(!withdrawPaused, "Withdraw is already paused");
        withdrawPaused = true;
        emit PauseWithdraw();
    }

    function unpauseWithdraw() public onlyRole(ADMIN_ROLE) {
        require(withdrawPaused, "Withdraw is not paused");
        withdrawPaused = false;
        emit UnpauseWithdraw();
    }
    //暂停和取消暂停 Claim
    function pauseClaim() public onlyRole(ADMIN_ROLE) {
        require(!claimPaused, "Claim is already paused");
        claimPaused = true;
        emit PauseClaim();
    }

    function unpauseClaim() public onlyRole(ADMIN_ROLE) {
        require(claimPaused, "Claim is not paused");
        claimPaused = false;
        emit UnpauseClaim();
    }

    //设置区块范围
    function setStartBlock(uint256 _startBlock) public onlyRole(ADMIN_ROLE) {
        require(_startBlock <= endBlock, "start block must be smaller than end block"));
        startBlock = _startBlock;
        emit SetStartBlock(_startBlock);
    }

    function setEndBlock(uint256 _endBlock) public onlyRole(ADMIN_ROLE) {
        require(startBlock <= _endBlock, "start block must be smaller than end block");
        endBlock = _endBlock;
        emit SetEndBlock(_endBlock);
    }

    //设置每区块的奖励（setRCCPerBlock）
    //功能：设置每个区块分发的 RCC 奖励数量。
    //权限：仅限管理员调用。
    //逻辑：
    //检查 _rccPerBlock 是否大于 0。
    //更新 rccPerBlock 值。
    //触发事件 SetRCCPerBlock。
    function setRccPerBlock(uint256 _rccPerBlock) public onlyRole(ADMIN_ROLE) {
        require(_rccPerBlock > 0, "invalid rccPerBlock");
        rccPerBlock = _rccPerBlock;
        emit setRccPerBlock(_rccPerBlock);
    }

//添加新质押池（addPool）
//功能：为质押奖励系统添加新的质押池。

//参数名称	数据类型	作用
//_stTokenAddress	address	指定该池子所用的质押代币地址。
//_poolWeight	uint256	设置池子的权重，用于奖励分配比例计算。
//_minDepositAmount	uint256	设置用户质押的最小数量限制。
//_unstakeLockedBlocks	uint256	质押后需锁仓的区块数量。
//_withUpdate	bool	是否在添加池子前更新其他池子的奖励状态。
//触发事件 AddPool。
    function addPool(address _stTokenAddress, uint256 _poolWeight, uint256 _minDepositAmount, uint256 _unstakeLockedBlocks, bool _withUpdate) public onlyRole(ADMIN_ROLE) {
        // 如果不是第一个池，质押代币地址必须有效（非 0 地址）。
        if (pool.length > 0) {
            require(_stTokenAddress != address(0x0), "invalid staking token address");
        } else {
            // 第一个池必须是 ETH 池，其地址设为 0。
            require(_stTokenAddress == address(0x0), "invalid staking token address");
        }

        // 确保锁仓的区块数量有效（大于 0）。
        require(_unstakeLockedBlocks > 0, "invalid withdraw locked blocks");

        // 确保当前区块号小于结束区块号。
        require(block.number < endBlock, "Already ended");

        // 如果需要更新池子状态，调用 `massUpdatePools`。
        if (_withUpdate) {
             massUpdatePools();
        }
        // 设置池子的最后奖励区块号：如果当前区块号大于起始区块号，则取当前区块号，否则取起始区块号。
        uint256 lastRewardBlock = block.number > startBlock ? block.number : startBlock;

        // 更新总池子权重。
        totalPoolWeight = totalPoolWeight + _poolWeight;

        // 创建新的池子对象并添加到 `pool` 数组。
        pool.push(Pool({
            stTokenAddress: _stTokenAddress,
            poolWeight: _poolWeight,
            lastRewardBlock: lastRewardBlock,
            accRCCPerST: 0, // 累计每个质押代币的 RCC 奖励初始化为 0。
            stTokenAmount: 0, // 初始质押代币数量为 0。
            minDepositAmount: _minDepositAmount,
            unstakeLockedBlocks: _unstakeLockedBlocks
        }));
        // 触发事件，记录添加的新池子信息。
        emit AddPool(_stTokenAddress, _poolWeight, lastRewardBlock, _minDepositAmount, _unstakeLockedBlocks);
    }

    /**
 * @notice 更新指定池子的相关信息（如最小质押数量和锁仓区块数量），仅管理员可调用。
 */
    function updatePool(uint256 _pid, uint256 _minDepositAmount, uint256 _unstakeLockedBlocks) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        pool[_pid].minDepositAmount = _minDepositAmount;
        pool[_pid].unstakeLockedBlocks = _unstakeLockedBlocks;

        emit UpdatePoolInfo(_pid, _minDepositAmount, _unstakeLockedBlocks);
    }

/**
 * @notice 更新池子的权重，仅管理员可调用。
 */
    function setPoolWeight(uint256 _pid, uint256 _poolWeight, bool _withUpdate) public onlyRole(ADMIN_ROLE) checkPid(_pid) {
        require(_poolWeight > 0, "invalid pool weight");
        if (_withUpdate) {
             massUpdatePools();
        }

        // 更新总池子权重，并修改指定池子的权重。
        totalPoolWeight = totalPoolWeight - pool[_poolWeight].poolWeight - _poolWeight;
        pool[_pid].poolWeight = _poolWeight;

        emit SetPoolWeight(_pid, _poolWeight, totalPoolWeight);
    }
    // ************************************** QUERY FUNCTION **************************************

    /**
    * @notice 返回池子数量（池子数组的长度）。
    */
    function poolLength() external view returns (uint256)  {
        return pool.length;
    }

    /**
 * @notice Return _from 到 _to 区块范围内所有区块奖励的总量。

 */
    function getMultiplier(uint256 _from, uint256 _to) public view returns (uint256 multiplier) {
        require(_from <= _to, "invalid block");

        // 如果 _from 小于奖励开始区块，则将 _from 调整为 startBlock
        if (_from < startBlock) {_from = startBlock;}

        // 如果 _to 大于奖励结束区块，则将 _to 调整为 endBlock
        if (_to > endBlock) {_to = endBlock;}

        // 再次确保调整后的 _from 小于或等于 _to
        require(_from <= _to, "end block must be greater than start block");

        // 使用安全的乘法计算区块范围内的奖励乘数，防止溢出
        bool success;
        //它表示 _from 到 _to 区块范围内所有区块奖励的总量。
        (success, multiplier) = (_to - _from).tryMul(rccPerBlock);
        require(success, "multiplier overflow");

    }

    //计算用户当前区块的待领取奖励
    function pendingRCC(uint256 _pid, address _user) external checkPid(_pid) view returns (uint256)  {
        return;
    }

    // 计算某用户在某个区块号时在指定池中的待发放奖励 RCC
    function pendingRCCByBlockNumber(uint256 _pid, address _user, uint256 _blockNumber) public checkPid(_pid) view returns (uint256) {
        Pool storage pool_ = pool[_pid];
        User storage user = user[_pid][_user];
        uint256 accRCCPerST = pool_.accRCCPerST; //每单位质押代币的累计 RCC 奖励（
        uint256 stSupply = pool_.stTokenAmount; //当前池中质押的代币总量

        if (_blockNumber > pool_.lastRewardBlock && stSupply != 0) {
            uint256 multiplier = getMultiplier(pool_.lastRewardBlock, _blockNumber); //计算池奖励的乘数
            uint256 rccForPool = multiplier * pool_.poolWeight / totalPoolWeight;  //计算池中的总奖励：
            accRCCPerST = accRCCPerST + rccForPool * (1 ether) / stSupply; //更新单位质押代币的累计奖励：
        }

        // 计算用户待发放奖励
        return user_.stAmount * accRCCPerST / (1 ether) - user_.finishedRCC + user_.pendingRCC;
    }

    /**
     * @notice 获取用户的下注量
     */
    function stakingBalance(uint256 _pid, address _user) external checkPid(_pid) view returns (uint256) {
        return user[_pid][_user].stAmount;
    }

    // 获取提现金额信息，包括未锁定金额和未锁定金额
    //查询用户在指定池中已申请的提取金额，以及当前可提取的金额（即已满足解锁区块的部分）。
    function withdrawAmount(uint256 _pid, address _user) public checkPid(_pid) view returns(uint256 requestAmount, uint256 pendingWithdrawAmount) {
        User storage user_ = user[_pid][_user];
        for (uint256 i = 0; i < user_.requests.length; i++) {
            //检查每个提取请求是否满足解锁条件。
            //unlockBlocks 提取请求对应的解锁区块号。 如果请求的解锁区块号小于等于当前区块号，则说明该请求的金额已满足解锁条件。
            if (user_.requests[i].unlockBlocks <= block.number) {
                pendingWithdrawAmount = pendingWithdrawAmount + user_.requests[i].amount;
            }
            requestAmount = requestAmount + user_.requests[i].amount;
        }

    }

    // ************************************** PUBLIC FUNCTION **************************************
    //更新指定池（_pid）的奖励变量，使其与当前区块的状态保持同步。
    //计算从上次更新奖励到当前区块产生的总奖励。
    //根据池的权重分配奖励，更新池的每股奖励（accRCCPerST）。
    //将池的最后奖励更新块号设置为当前块号。
    function updatePool(uint256 _pid) public checkPid(_pid)  {
         Pool storage pool_ = pool[_pid];
        if(block.number <= pool[_pid]){
            return;
        }
            // 区块范围内所有区块奖励的总量。

        (bool success1, uint256 totalRCC) = getMultiplier(pool_.lastRewardBlock, block.number).tryMul(pool_.poolWeight);
        require(success1, "overflow");

       // 分配奖励（按总权重比例）
        (success1, totalRCC) = totalRCC.tryDiv(totalPoolWeight);
        require(success1, "overflow");

        //. 获取池的总质押量
        uint256 stSupply = pool_.stTokenAmount; //当前池中质押的代币总量

        // 计算并更新每股奖励
        if (stSupply > 0) {
            (bool success2, uint256 totalRCC_) = totalRCC.tryMul(1 ether);
            require(success2, "overflow");

            //将奖励按池的总质押量分配，计算每单位质押代币的奖励
            (success2, totalRCC_) = totalRCC_.tryDiv(stSupply);
            require(success2, "overflow");

            //将计算出的每单位奖励累加到池的累积每股奖励
            (bool success3, uint256 accRCCPerST) = pool_.accRCCPerST.tryAdd(totalRCC_);
            require(success3, "overflow");
            pool_.accRCCPerST = accRCCPerST;



        }
//pool_.lastRewardBlock = block.number;
        pool_.lastRewardBlock = block.number;
        emit UpdatePool(_pid, pool_.lastRewardBlock, totalRCC);
    }

    function massUpdatePools()  {
        uint256 length = pool.length;
        for (uint256 pid = 0; pid < length; ++pid) {
            updatePool(pid);
        }
    }

    function depositETH() public whenNotPaused() payable {
        Pool storage pool_ = pool[ETH_PID];
        require(pool_.stTokenAddress == address(0x0), "invalid staking token address");

        //msg.value 是用户通过交易发送的 ETH 金额（以 wei 为单位）。
        //将发送的 ETH 金额存入变量 _amount 中。
        uint256 _amount = msg.value;
        require(_amount >= pool_.minDepositAmount, "deposit amount is too small");

        _deposit(ETH_PID, _amount);
    }


    //用户向指定池（pool）中存入 ERC-20 代币
    function deposit(uint256 _pid, uint256 _amount) public whenNotPaused() checkPid(_pid) {
        require(_pid != 0, "deposit not support ETH staking");
        Pool storage pool_ = pool[_pid];
        require(_amount > pool_.minDepositAmount, "deposit amount is too small");

        if(_amount > 0) {
            IERC20(pool_.stTokenAddress).safeTransferFrom(address(msg.sender), address(this), _amount);
        }

    }

    }