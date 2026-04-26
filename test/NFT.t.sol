// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/NFT.sol";

contract NFTTest is Test {
    NFT nft;

    address owner = address(1);
    address user = address(2);

    function setUp() public {
        vm.prank(owner);
        nft = new NFT();
    }

    /* -------------------------- Mint Tests -------------------------- */
    
    function testMintSuccess() public {
        vm.deal(user, 1 ether);

        vm.prank(user);
        uint256 tokenId = nft.mint{value: 0.01 ether}("ipfs://token1");

        assertEq(tokenId, 1);
        assertEq(nft.ownerOf(1), user);
        assertEq(nft.tokenURI(1), "ipfs://token1");
    }

    function testMintFailsInsufficientPayment() public {
        vm.deal(user, 1 ether);

        vm.prank(user);
        vm.expectRevert("Insufficient payment");
        nft.mint{value: 0.005 ether}("ipfs://token1");
    }

    function testMintFailsWhenMaxSupplyReached() public {
        vm.deal(user, 200 ether);

        vm.startPrank(user);

        for (uint256 i = 0; i < nft.MAX_SUPPLY(); i++) {
            nft.mint{value: 0.01 ether}("uri");
        }

        vm.expectRevert("Max supply reached");
        nft.mint{value: 0.01 ether}("uri");

        vm.stopPrank();
    }

    function testMintIncrementsTokenId() public {
        vm.deal(user, 1 ether);

        vm.startPrank(user);
        uint256 tokenId1 = nft.mint{value: 0.01 ether}("uri1");
        uint256 tokenId2 = nft.mint{value: 0.01 ether}("uri2");
        vm.stopPrank();

        assertEq(tokenId2, tokenId1 + 1);
    }

    /* ---------------------- Mint Price Tests ----------------------- */

    function testSetMintPrice() public {
        vm.prank(owner);
        nft.setMintPrice(0.02 ether);

        assertEq(nft.mintPrice(), 0.02 ether);
    }

    function testSetMintPriceFailsIfNotOwner() public {
        vm.prank(user);
        vm.expectRevert();
        nft.setMintPrice(0.02 ether);
    }

    function testSetMintPriceFailsIfZero() public {
        vm.prank(owner);
        vm.expectRevert("Price must be greater than zero");
        nft.setMintPrice(0);
    }

    /* ------------------------ Withdraw Tests ------------------------ */

    function testWithdraw() public {
        vm.deal(user, 1 ether);

        vm.prank(user);
        nft.mint{value: 0.01 ether}("uri");

        uint256 ownerBalanceBefore = owner.balance;

        vm.prank(owner);
        nft.withdraw();

        assertEq(address(nft).balance, 0);
        assertEq(owner.balance, ownerBalanceBefore + 0.01 ether);
    }

    function testWithdrawFailsIfNoFunds() public {
        vm.prank(owner);
        vm.expectRevert("No funds to withdraw");
        nft.withdraw();
    }

    function testWithdrawFailsIfNotOwner() public {
        vm.deal(user, 1 ether);

        vm.prank(user);
        nft.mint{value: 0.01 ether}("uri");

        vm.prank(user);
        vm.expectRevert();
        nft.withdraw();
    }
}