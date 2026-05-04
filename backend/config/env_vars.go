package config

import "github.com/ethereum/go-ethereum/common"

const AuctionAddressHex = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

var AuctionAddress = common.HexToAddress(AuctionAddressHex)

const MarketplaceHex = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"

var MarketplaceAddress = common.HexToAddress(MarketplaceHex)

const EthNodeWsUrl = "ws://127.0.0.1:8545/"
