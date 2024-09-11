 
const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");
module.exports = buildModule("ShoppingModule", (m) => {
    const shopping = m.contract("Shopping",[]);
    m.call(shopping, "Status", []);
    return {shopping};
});