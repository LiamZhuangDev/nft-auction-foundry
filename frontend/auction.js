import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { getContracts } from "./contracts.js";
import { API_URL } from "./config.js";

export async function fetchAuctions() {
  try {
    const res = await fetch(`${API_URL}/auctions`);

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    const result = await res.json();
    const auctions = result.data;

    console.log("Auctions:", auctions);
    
    renderAuctions(auctions);
  } catch (err) {
    console.error("Error fetching auctions:", err);
  }
}

function renderAuctions(auctions) {
  const container = document.getElementById("auctions");
  container.innerHTML = "";

  auctions.forEach(a => {
    const now = Math.floor(Date.now() / 1000);
    const secondsRemaining = Number(a.EndTime) - now;
    const status = a.Active && secondsRemaining > 0
      ? `Active (${secondsRemaining}s left)`
      : "Ended";

    const div = document.createElement("div");

    const finalPrice = a.FinalPrice && a.FinalPrice !== ""
      ? ethers.formatEther(a.FinalPrice)
      : "N/A";

    div.innerHTML = `
      <p><strong>Auction ID: </strong> ${a.AuctionID}</p>
      <p><strong>Seller: </strong> ${a.Seller}</p>
      <p><strong>NFT: </strong> ${a.NftContract}</p>
      <p><strong>Token ID: </strong> ${a.TokenId}</p>
      <p><strong>Start Price: </strong> ${ethers.formatEther(a.StartPrice)} ETH</p>
      <p><strong>Final Price: </strong> ${finalPrice} ETH</p>
      <p><strong>End Time: </strong> ${new Date(a.EndTime * 1000).toLocaleString()}</p>
      <p><strong>Status: </strong> ${status}</p>
    `;
    container.appendChild(div);
  })
}

export async function fetchBids() {
  try {
    const auctionId = document.getElementById("auctionIdInput").value;
    if (!auctionId) {
      alert("Please enter an auction ID");
      return;
    }

    const res = await fetch(`${API_URL}/auctions/${auctionId}/bids`);

    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }

    const result = await res.json();
    const bids = result.data;

    console.log("Bids:", bids);
    
    renderBids(bids);
  } catch (err) {
    console.error("Error fetching bids:", err);
  }
}

function renderBids(bids) {
  const container = document.getElementById("bids");
  container.innerHTML = "";

  bids.forEach(b => {
    const div = document.createElement("div");
    div.innerHTML = `
      <p><strong>Auction ID: </strong> ${b.AuctionID}</p>
      <p><strong>Bidder: </strong> ${b.Bidder}</p>
      <p><strong>Amount: </strong> ${ethers.formatEther(b.Amount)} ETH</p>
      <p><strong>Timestamp: </strong> ${new Date(b.Timestamp * 1000).toLocaleString()}</p>
    `;
    container.appendChild(div);
  })
}

export async function placeBid() {
  const { auction } = getContracts();

  const auctionId = document.getElementById("auctionId").value;
  const bidAmount = document.getElementById("bidAmount").value;
  const auctionData = await auction.auctions(auctionId);
  const endTime = Number(auctionData.endTime);
  const now = Math.floor(Date.now() / 1000);

  if (now >= endTime) {
    throw new Error(`Auction ${auctionId} ended at ${new Date(endTime * 1000).toLocaleString()}`);
  }

  const tx = await auction.placeBid(auctionId, {
    value: ethers.parseEther(bidAmount)
  });

  await tx.wait();
  console.log("Bid Placed");
}

export async function finalizeAuction() {
  const { auction } = getContracts();

  const auctionId = document.getElementById("finalizeAuctionId").value;

  const tx = await auction.finalizeAuction(auctionId);
  await tx.wait();

  console.log("Auction Finalized");
}
