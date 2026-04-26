// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import "openzeppelin-contracts/contracts/interfaces/IERC721.sol";
import "openzeppelin-contracts/contracts/interfaces/IERC721Receiver.sol";

contract NFTAuctionHouse is IERC721Receiver {
    struct Auction {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 highestBid;
        address highestBidder;
        uint256 endTime;
        bool active;
    }

    Auction[] public auctions;
    uint256 public constant AuctionFee = 500; // 5% = 500 / 10,000
    address public feeRecipient;

    mapping(uint256 => mapping(address => uint256)) public pendingReturns; // auctionId => bidder => amount
    mapping(address => mapping(uint256 => bool)) public activeAuctions; // nftContract => tokenId => isActive

    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 end);
    
    event BidPlaced(
        uint256 indexed auctionId, 
        address indexed bidder, 
        uint256 amount);

    event AuctionEnded(
        uint256 indexed auctionId, 
        address indexed winner, 
        uint256 amount);

    event Withdrawal(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount);

    constructor(address _feeRecipient) {
        feeRecipient = _feeRecipient;
    }

    function createAuction(
        address seller,
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 duration
    ) external returns (uint256 auctionId)
    {
        auctions.push(Auction({
            seller: seller,
            nftContract: nftContract,
            tokenId: tokenId,
            highestBid: startPrice,
            highestBidder: address(0),
            endTime: block.timestamp + duration,
            active: true
        }));
        auctionId = auctions.length - 1;
        activeAuctions[nftContract][tokenId] = true;

        emit AuctionCreated(auctionId, seller, nftContract, tokenId, startPrice, block.timestamp + duration);
    }
    
    function placeBid(uint256 auctionId) external payable {
        require(auctionId < auctions.length, "Invalid auction ID");
        Auction storage a = auctions[auctionId];

        require(a.active, "Auction is not active");
        require(block.timestamp < a.endTime, "Auction has ended");
        require(msg.value > a.highestBid, "Bid must be higher than current highest");

        // Refund previous highest bidder
        if (a.highestBidder != address(0)) {
            pendingReturns[auctionId][a.highestBidder] += a.highestBid;
        }

        // Update highest bid and bidder
        a.highestBidder = msg.sender;
        a.highestBid = msg.value;

        emit BidPlaced(auctionId, msg.sender, msg.value);
    }

    // pull over push pattern for withdrawals
    function withdrawBid(uint256 auctionId) external {
        require(auctionId < auctions.length, "Invalid auction ID");
        require(block.timestamp >= auctions[auctionId].endTime, "Auction is still active");

        // CEI pattern: Check, Effects, Interactions
        uint256 amount = pendingReturns[auctionId][msg.sender];
        require(amount > 0, "No funds to withdraw");

        pendingReturns[auctionId][msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Withdrawal failed");

        emit Withdrawal(auctionId, msg.sender, amount);
    }

    function finalizeAuction(uint256 auctionId) external {
        require(auctionId < auctions.length, "Invalid auction Id");
        Auction storage a = auctions[auctionId];

        require(a.active, "Auction is not active");
        require(block.timestamp >= a.endTime, "Auction has not ended");
        require(msg.sender == a.seller || msg.sender == a.highestBidder, "Only seller or highest bidder can finalize the auction");

        a.active = false;
        activeAuctions[a.nftContract][a.tokenId] = false;

        if (a.highestBidder != address(0)) {
            // Calculate fee and seller proceeds
            uint256 fee = (a.highestBid * AuctionFee) / 10000;
            uint256 sellerProceeds = a.highestBid - fee;

            // Transfer NFT to highest bidder
            IERC721(a.nftContract).safeTransferFrom(address(this), a.highestBidder, a.tokenId);

            // Transfer funds to seller and fee recipient
            (bool success, ) = a.seller.call{value: sellerProceeds}("");
            require(success, "Payment to seller failed");

            (bool success2, ) = feeRecipient.call{value: fee}("");
            require(success2, "Payment to fee recipient failed");
        } else {
            // No bids were placed, return NFT to seller
            IERC721(a.nftContract).safeTransferFrom(address(this), a.seller, a.tokenId);
        }

        emit AuctionEnded(auctionId, a.highestBidder, a.highestBid);
    }

    function isActive(address nftContract, uint256 tokenId) external view returns (bool) {
        require(nftContract != address(0), "Invalid NFT contract");
        require(tokenId > 0, "Invalid token ID");
        return activeAuctions[nftContract][tokenId];
    }

    function getAuctionCount() external view returns (uint256) {
        return auctions.length;
    }

    // onERC721Received is ONLY required if the receiver is a CONTRACT
    function onERC721Received(
        address /*operator*/, 
        address /*from*/, 
        uint256 /*tokenId*/, 
        bytes calldata /*data*/) external pure override returns (bytes4) {
            return IERC721Receiver.onERC721Received.selector;
    }
}