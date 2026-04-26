// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "openzeppelin-contracts/contracts/interfaces/IERC721.sol";

interface IAuctionHouse {
    function createAuction(
        address seller,
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 duration
    ) external returns (uint256 auctionId);

    function isActive(address nftContract, uint256 tokenId) external view returns (bool);
}

// Design Overview:
// User
//  ↓
// Marketplace (entry point / orchestrator)
//  └─ Delegates auctions → AuctionHouse
//                           ├─ createAuction
//                           ├─ placeBid
//                           └─ endAuction
contract NFTMarketplace {
    struct Listing {
        address seller;
        address nftContract;
        uint256 tokenId;
        bool isListing;
    }

    Listing[] public listings;
    mapping(address => mapping(uint256 => bool)) public activeListings; // nftContract => tokenId => isActive

    IAuctionHouse public auctionHouse;

    event Listed(address indexed seller, address indexed nftContract, uint256 tokenId, uint256 listingId);
    event Delisted(address indexed seller, address indexed nftContract, uint256 tokenId, uint256 listingId);

    constructor(address _auctionContract) {
        // Type cast, treat this address as if it implements the IAuctionHouse interface
        // Solidity does zero runtime checking, so if the contract doesn't support functions in interface, transaction reverts.
        auctionHouse = IAuctionHouse(_auctionContract); 
    }

    function listNFT(address nftContract, uint256 tokenId)
        external returns (uint256)
    {
        require(tokenId > 0, "Invalid token ID");
        require(nftContract != address(0), "Invalid NFT contract");
        require(!activeListings[nftContract][tokenId], "NFT is already listed");

        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Only owner can list NFT");

        listings.push(Listing({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            isListing: true
        }));

        activeListings[nftContract][tokenId] = true;

        uint256 listingId = listings.length - 1;
        emit Listed(msg.sender, nftContract, tokenId, listingId);

        return listingId;
    }

    function delistNFT(uint256 listingId) external {
        require(listingId < listings.length, "Invalid listing ID");

        Listing storage l = listings[listingId];
        require(l.isListing, "NFT is not listing");
        require(l.seller == msg.sender, "Only seller can delist NFT");

        l.isListing = false;
        activeListings[l.nftContract][l.tokenId] = false;

        emit Delisted(msg.sender, l.nftContract, l.tokenId, listingId);
    }

    function createAuction(
        uint256 listingId, 
        uint256 startPrice, 
        uint256 duration) 
        external returns (uint256 auctionId)
    {
        // Validate inputs
        require(listingId < listings.length, "Invalid listing ID");
        require(startPrice > 0, "Start price must be greater than zero");
        require(duration > 0, "Duration must be greater than zero");

        Listing storage l = listings[listingId];
        address nftContract = l.nftContract;
        uint256 tokenId = l.tokenId;
        require(!auctionHouse.isActive(nftContract, tokenId),"NFT already in auction");

        // verify ownership and approval
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Only owner can create auction");
        require(nft.getApproved(tokenId) == address(this) || nft.isApprovedForAll(msg.sender, address(this)), "Marketplace must be approved to transfer NFT");
        
        // transfer NFT and create auction
        nft.safeTransferFrom(msg.sender, address(auctionHouse), tokenId);
        auctionId = auctionHouse.createAuction(msg.sender, nftContract, tokenId, startPrice, duration);
    }

    function getListingCount() external view returns (uint256) {
        return listings.length;
    }
}
