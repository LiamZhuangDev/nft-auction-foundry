package config

import "github.com/ethereum/go-ethereum/common"

const AuctionAddressHex = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"

var AuctionAddress = common.HexToAddress(AuctionAddressHex)

const MarketplaceHex = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

var MarketplaceAddress = common.HexToAddress(MarketplaceHex)

const EthNodeWsUrl = "ws://127.0.0.1:8545/"
