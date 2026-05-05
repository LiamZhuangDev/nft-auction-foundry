import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";

export let provider;
export let signer;

export async function connectWallet() {
  if (!window.ethereum) {
    alert("Install Wallet Extension (e.g. MetaMask) to use this app");
    return;
  }

  provider = new ethers.BrowserProvider(window.ethereum);

  await provider.send("eth_requestAccounts", []);
  signer = await provider.getSigner();

  console.log("Connected:", await signer.getAddress());
}