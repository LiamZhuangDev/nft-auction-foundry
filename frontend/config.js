export const NFT_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
export const MARKETPLACE_ADDRESS = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";
export const AUCTION_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";

export const nftABI = [
    "function mint(string uri) payable returns (uint256)",
    "function approve(address to, uint256 tokenId)",
];
export const marketplaceABI = [
    "function listNFT(address nftContract, uint256 tokenId) returns (uint256)",
    "function createAuction(uint256 listingId, uint256 startPrice, uint256 duration) returns (uint256)"
];
export const auctionABI = [
    "function auctions(uint256 auctionId) view returns (address seller, address nftContract, uint256 tokenId, uint256 highestBid, address highestBidder, uint256 endTime, bool active)",
    "function placeBid(uint256 auctionId) payable",
    "function finalizeAuction(uint256 auctionId)",
    "function withdrawBid(uint256 auctionId)",
];

export const API_URL = "http://localhost:8081";
