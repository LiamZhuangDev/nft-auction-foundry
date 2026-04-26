### Key Foundry Test concepts
- vm.prank(addr) → simulate caller, applies to *ONE* single call
- vm.startPrank(addr)/vm.stopPrank(), applies to *ALL* calls until stopped
- vm.deal(addr, amount) → give ETH
- vm.expectRevert() → expect failure
- assertEq() → assertions
