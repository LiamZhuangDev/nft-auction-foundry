// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Script.sol";
import "forge-std/console2.sol";

import "../src/NFT.sol";
import "../src/NFTAuctionHouse.sol";
import "../src/NFTMarketplace.sol";

contract Deploy is Script {
    function run() external returns (NFT nft, NFTAuctionHouse auctionHouse, NFTMarketplace marketplace) {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address feeRecipient = vm.envOr("FEE_RECIPIENT", vm.addr(deployerPrivateKey));

        vm.startBroadcast(deployerPrivateKey);

        nft = new NFT();
        auctionHouse = new NFTAuctionHouse(feeRecipient);
        marketplace = new NFTMarketplace(address(auctionHouse));

        vm.stopBroadcast();

        console2.log("NFT deployed at:", address(nft));
        console2.log("NFTAuctionHouse deployed at:", address(auctionHouse));
        console2.log("NFTMarketplace deployed at:", address(marketplace));
        console2.log("Fee recipient:", feeRecipient);
    }
}
