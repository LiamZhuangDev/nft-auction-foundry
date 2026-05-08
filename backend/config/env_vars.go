package config

import "github.com/ethereum/go-ethereum/common"

const AuctionAddressHex = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512" // Anvil local network
// const AuctionAddressHex = "0xc2214d88C9ae33DfC275F088a5808b321AF43972" // Anvil Mainnet Fork, changes on each deployment
var AuctionAddress = common.HexToAddress(AuctionAddressHex)

const MarketplaceHex = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0" // Anvil local network
// const MarketplaceHex = "0xABd78942Eca28f15E0F37E6Bae7DA6879fD1257a" // Anvil Mainnet Fork, changes on each deployment
var MarketplaceAddress = common.HexToAddress(MarketplaceHex)

const EthNodeWsUrl = "ws://127.0.0.1:8545/"

const ChainlinkFeedAddressHex = "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419" // Chainlink Ethereum Mainnet Address
var ChainlinkFeedAddress = common.HexToAddress(ChainlinkFeedAddressHex)
