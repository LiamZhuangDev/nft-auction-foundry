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


