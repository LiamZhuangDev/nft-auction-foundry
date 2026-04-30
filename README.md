### Start Anvil node in terminal 1
```bash
anvil
```
### Deploy contracts to the running node in terminal 2
```bash
export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # Anvil’s default first account

forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast
```

### Foundry Test functions
- Must start with `test`
```solidity
function testMint() public {}
function testWithdraw() public {}
```
- Fuzz tests, Foundry generates random inputs and run many times
```solidity
function testFuzz_Mint(uint256 amount) public {}
```

### Key Foundry Test concepts
- vm.prank(addr) → simulate caller, applies to *ONE* single call
- vm.startPrank(addr)/vm.stopPrank(), applies to *ALL* calls until stopped
- vm.deal(addr, amount) → give ETH
- vm.expectRevert() → expect failure
- assertEq() → assertions

### Foundry built-in hook
- setUp(), it runs before every single test
```solidity
function setUp() public {}
```


