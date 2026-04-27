// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/NFTMarketplace.sol";
import "../src/NFT.sol";
import "openzeppelin-contracts/contracts/token/ERC721/IERC721Receiver.sol";

/* ------------------------------------- Mock AuctionHouse ----------------------------------------- */

contract MockAuctionHouse is IAuctionHouse, IERC721Receiver {
    uint256 public nextId;
    mapping(address => mapping(uint256 => bool)) public active;

    function createAuction(address, address nftContract, uint256 tokenId, uint256, uint256)
        external
        returns (uint256 auctionId)
    {
        auctionId = nextId++;
        active[nftContract][tokenId] = true;
    }

    function isActive(address nftContract, uint256 tokenId) external view returns (bool) {
        return active[nftContract][tokenId];
    }

    function onERC721Received(address, address, uint256, bytes calldata) external pure override returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }
}

/* ------------------------------------- Tests ----------------------------------------- */

contract NFTMarketplaceTest is Test {
    NFT nft;
    uint256 tokenId;
    NFTMarketplace marketplace;
    MockAuctionHouse auctionHouse;

    address user = address(1);
    address other = address(2);

    function setUp() public {
        nft = new NFT();
        auctionHouse = new MockAuctionHouse();
        marketplace = new NFTMarketplace(address(auctionHouse));

        // mint NFT to user
        vm.deal(user, 10 ether);
        vm.prank(user);
        tokenId = nft.mint{value: 0.01 ether}("uri"); // tokenId = 1
    }

    /* ---------------------- Listing Tests---------------------- */

    function testListNFT() public {
        vm.prank(user);
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        (address seller, address nftAddr, uint256 id, bool isListing) = marketplace.listings(listingId);

        assertEq(seller, user);
        assertEq(nftAddr, address(nft));
        assertEq(id, tokenId);
        assertTrue(isListing);
        assertTrue(marketplace.activeListings(nftAddr, id));
    }

    function testListFailsIfInvalidTokenId() public {
        vm.prank(user);
        vm.expectRevert("Invalid token ID");
        marketplace.listNFT(address(nft), 0);
    }

    function testListFailsIfNotOwner() public {
        vm.prank(other);
        vm.expectRevert("Only owner can list NFT");
        marketplace.listNFT(address(nft), tokenId);
    }

    function testListFailsIfAlreadyListed() public {
        vm.startPrank(user);
        marketplace.listNFT(address(nft), tokenId);

        vm.expectRevert("NFT is already listed");
        marketplace.listNFT(address(nft), tokenId);
        vm.stopPrank();
    }

    /* ---------------------- Delist Tests---------------------- */
    function testDelistNFT() public {
        vm.prank(user);
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        vm.prank(user);
        marketplace.delistNFT(listingId);

        (,,, bool isListing) = marketplace.listings(listingId);

        assertFalse(isListing);
        assertFalse(marketplace.activeListings(address(nft), tokenId));
    }

    function testDelistFailsIfInvalidId() public {
        vm.expectRevert("Invalid listing ID");
        marketplace.delistNFT(999);
    }

    function testDelistFailsIfNotListed() public {
        vm.prank(user);
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        vm.prank(user);
        marketplace.delistNFT(listingId);

        vm.expectRevert("NFT is not listing");
        marketplace.delistNFT(listingId);
    }

    function testDelistFailsIfNotOwner() public {
        vm.prank(user);
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        vm.prank(other);
        vm.expectRevert("Only seller can delist NFT");
        marketplace.delistNFT(listingId);
    }

    /* -------------------- Create Auction Tests -------------------- */

    function testCreateAuction() public {
        vm.startPrank(user);

        // listing
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);
        // approve marketplace
        nft.approve(address(marketplace), tokenId);
        // create auction
        uint256 auctionId = marketplace.createAuction(listingId, 1 ether, 1 days);

        assertEq(auctionId, 0);
        assertTrue(auctionHouse.isActive(address(nft), tokenId));

        vm.stopPrank();
    }

    function testCreateAuctionFailsIfNotOwner() public {
        vm.prank(user);
        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        vm.prank(other);
        vm.expectRevert("Only owner can create auction");
        marketplace.createAuction(listingId, 1 ether, 1 days);
    }

    function testCreateAuctionFailsWithoutApproval() public {
        vm.startPrank(user);

        uint256 listingId = marketplace.listNFT(address(nft), tokenId);

        vm.expectRevert("Marketplace must be approved to transfer NFT");
        marketplace.createAuction(listingId, 1 ether, 1 days);

        vm.stopPrank();
    }

    function testCreateAuctionFailsIfAlreadyInAuction() public {
        vm.startPrank(user);

        uint256 listingId = marketplace.listNFT(address(nft), tokenId);
        nft.approve(address(marketplace), tokenId);
        marketplace.createAuction(listingId, 1 ether, 1 days);

        vm.expectRevert("NFT already in auction");
        marketplace.createAuction(listingId, 1 ether, 1 days);

        vm.stopPrank();
    }

    function testCreateAuctionFailsInvalidParams() public {
        vm.startPrank(user);

        uint256 listingId = marketplace.listNFT(address(nft), tokenId);
        nft.approve(address(marketplace), tokenId);

        vm.expectRevert("Invalid listing ID");
        marketplace.createAuction(999, 1 ether, 1 days);

        vm.expectRevert("Start price must be greater than zero");
        marketplace.createAuction(listingId, 0, 1 days);

        vm.expectRevert("Duration must be greater than zero");
        marketplace.createAuction(listingId, 1 ether, 0);
    }

    /* -------------------- View Tests -------------------- */
    function testGetListingCount() public {
        vm.prank(user);
        marketplace.listNFT(address(nft), tokenId);

        assertEq(marketplace.getListingCount(), 1);
    }
}
