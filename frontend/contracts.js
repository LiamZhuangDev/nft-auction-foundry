import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { signer } from "./wallet.js";
import {
  NFT_ADDRESS,
  MARKETPLACE_ADDRESS,
  AUCTION_ADDRESS,
  nftABI,
  marketplaceABI,
  auctionABI
} from "./config.js";

export function getContracts() {
  if (!signer) throw new Error("Wallet not connected");

  return {
    nft: new ethers.Contract(NFT_ADDRESS, nftABI, signer),
    marketplace: new ethers.Contract(MARKETPLACE_ADDRESS, marketplaceABI, signer),
    auction: new ethers.Contract(AUCTION_ADDRESS, auctionABI, signer)
  };
}