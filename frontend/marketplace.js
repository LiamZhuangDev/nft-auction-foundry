import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { getContracts } from "./contracts.js";
import { NFT_ADDRESS, API_URL } from "./config.js";

export async function fetchListings() {
    try {
        const res = await fetch(`${API_URL}/listings`);

        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }

        const result = await res.json();
        const listings = result.data;

        console.log("Listings:", listings);

        renderListings(listings);
    } catch (err) {
        console.error("Error fetching listings:", err);
    }
}

function renderListings(listings) {
    const container = document.getElementById("listings");
    container.innerHTML = "";

    listings.forEach(l => {
        const div = document.createElement("div");
        div.innerHTML = `
            <p><strong>Listing ID: </strong> ${l.ID}</p>
            <p><strong>Seller: </strong> ${l.Seller}</p>
            <p><strong>NFT: </strong> ${l.NftContract}</p>
            <p><strong>Token ID: </strong> ${l.TokenId}</p>
            <p><strong>Status: </strong> ${l.Active ? "Active" : "Inactive"}</p>
        `;
        container.appendChild(div);
    })
}

export async function listNFT() {
    const { marketplace } = getContracts();

    const tokenId = document.getElementById("listTokenId").value;

    const tx = await marketplace.listNFT(NFT_ADDRESS, tokenId);
    await tx.wait();

    console.log("NFT listed for sale");
}

export async function createAuction() {
    const { marketplace } = getContracts();

    const listingId = document.getElementById("listingId").value;
    const startPrice = document.getElementById("startPrice").value;
    const duration = document.getElementById("duration").value;

    const tx = await marketplace.createAuction(
        listingId,
        ethers.parseEther(startPrice),
        duration
    );
    await tx.wait();

    console.log("Auction created");
}