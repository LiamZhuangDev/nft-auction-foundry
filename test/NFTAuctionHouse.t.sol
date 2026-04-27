// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/NFT.sol";
import "../src/NFTAuctionHouse.sol";

contract NFTAuctionHouseTest is Test {
    NFT nft;
    uint256 tokenId;
    NFTAuctionHouse auctionHouse;

    address seller = address(1);
    address bidder1 = address(2);
    address bidder2 = address(3);
    address feeRecipient = address(4);

    function setUp() public {
        nft = new NFT();
        auctionHouse = new NFTAuctionHouse(feeRecipient);

        vm.deal(seller, 10 ether);
        vm.deal(bidder1, 10 ether);
        vm.deal(bidder2, 10 ether);

        // mint NFT to seller
        vm.prank(seller);
        tokenId = nft.mint{value: 0.01 ether}("uri");

        // approve auction house
        vm.prank(seller);
        nft.approve(address(auctionHouse), tokenId);

        // AuctionHouse should assume NFT is already transferred in
        vm.prank(seller);
        nft.safeTransferFrom(seller, address(auctionHouse), tokenId);
    }

    /* -------------- Create Auction Test -------------- */
    function testCreateAuction() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        assertEq(auctionId, 0);
        assertTrue(auctionHouse.activeAuctions(address(nft), tokenId));
    }

    /* ------------------ Bid Tests ------------------ */
    function testPlaceBid() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        auctionHouse.placeBid{value: 2 ether}(auctionId);

        (,,, uint256 highestBid, address highestBidder,,) = auctionHouse.auctions(auctionId);

        assertEq(highestBid, 2 ether);
        assertEq(highestBidder, bidder1);
    }

    function testPlaceBidFailsIfValueTooLow() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        vm.expectRevert("Bid must be higher than current highest");
        auctionHouse.placeBid{value: 1 ether}(auctionId);
    }

    function testBidFailsIfAuctionEnded() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), 1, 1 ether, 1 days);

        vm.warp(block.timestamp + 2 days);

        vm.prank(bidder1);
        vm.expectRevert("Auction has ended");
        auctionHouse.placeBid{value: 2 ether}(auctionId);
    }

    function testRefundPreviousBidder() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        auctionHouse.placeBid{value: 2 ether}(auctionId);

        vm.prank(bidder2);
        auctionHouse.placeBid{value: 3 ether}(auctionId);

        assertEq(auctionHouse.pendingReturns(auctionId, bidder1), 2 ether);
    }

    /* ------------------ Withdraw Tests ------------------ */
    function testWithdrawBid() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        auctionHouse.placeBid{value: 2 ether}(auctionId);

        vm.prank(bidder2);
        auctionHouse.placeBid{value: 3 ether}(auctionId);

        vm.warp(block.timestamp + 2 days);

        uint256 before = bidder1.balance;

        vm.prank(bidder1);
        auctionHouse.withdrawBid(auctionId);

        assertEq(bidder1.balance, before + 2 ether);
        assertEq(auctionHouse.pendingReturns(auctionId, bidder1), 0);
    }

    function testWithdrawFailsIfActive() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        auctionHouse.placeBid{value: 2 ether}(auctionId);

        vm.prank(bidder1);
        vm.expectRevert("Auction is still active");
        auctionHouse.withdrawBid(auctionId);
    }

    /* ------------------ Finalize Auction ------------------ */

    function testFinalizeWithBids() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.prank(bidder1);
        auctionHouse.placeBid{value: 2 ether}(auctionId);

        vm.warp(block.timestamp + 2 days);

        uint256 sellerBefore = seller.balance;
        uint256 feeBefore = feeRecipient.balance;

        vm.prank(seller);
        auctionHouse.finalizeAuction(auctionId);

        // NFT transferred
        assertEq(nft.ownerOf(tokenId), bidder1);

        uint256 fee = (2 ether * auctionHouse.AuctionFee()) / 10000;
        uint256 sellerProceeds = 2 ether - fee;

        assertEq(seller.balance, sellerBefore + sellerProceeds);
        assertEq(feeRecipient.balance, feeBefore + fee);
    }

    function testFinalizeNoBids() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.warp(block.timestamp + 2 days);

        vm.prank(seller);
        auctionHouse.finalizeAuction(auctionId);

        // NFT returned to seller
        assertEq(nft.ownerOf(tokenId), seller);
    }

    function testFinalizeFailsIfNotAllowed() public {
        vm.prank(seller);
        uint256 auctionId = auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        vm.warp(block.timestamp + 2 days);

        vm.prank(bidder1);
        vm.expectRevert("Only seller or highest bidder can finalize the auction");
        auctionHouse.finalizeAuction(auctionId);
    }

    /* ------------------ isActive ------------------ */

    function testIsActive() public {
        vm.prank(seller);
        auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        assertTrue(auctionHouse.isActive(address(nft), tokenId));
    }

    function testIsActiveFailsInvalid() public {
        vm.expectRevert("Invalid NFT contract");
        auctionHouse.isActive(address(0), 1);

        vm.expectRevert("Invalid token ID");
        auctionHouse.isActive(address(nft), 0);
    }

    /* ------------------ Count ------------------ */

    function testGetAuctionCount() public {
        vm.prank(seller);
        auctionHouse.createAuction(seller, address(nft), tokenId, 1 ether, 1 days);

        assertEq(auctionHouse.getAuctionCount(), 1);
    }
}
