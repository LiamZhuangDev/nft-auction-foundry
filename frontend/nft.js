import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { getContracts } from "./contracts.js";
import { MARKETPLACE_ADDRESS } from "./config.js";

export async function mintNFT() {
    const { nft } = getContracts();

    const uri = document.getElementById("tokenURI").value;

    const tx = await nft.mint(uri, {
        value: ethers.parseEther("0.01") // Minting fee
    });
    await tx.wait();
    console.log("NFT minted");
}

export async function approveNFT() {
    const { nft } = getContracts();

    const tokenId = document.getElementById("listTokenId").value;

    const tx = await nft.approve(MARKETPLACE_ADDRESS, tokenId);
    await tx.wait();

    console.log("NFT approved for marketplace");
}