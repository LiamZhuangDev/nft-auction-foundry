package config

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// var AuctionCreatedSig = crypto.Keccak256Hash([]byte("AuctionCreated(uint256,address,address,uint256,uint256,uint256)"))
// var BidPlacedSig = crypto.Keccak256Hash([]byte("BidPlaced(uint256,address,uint256)"))
// var AuctionFinalizeSig = crypto.Keccak256Hash([]byte("AuctionEnded(uint256,address,uint256)"))

type PollerConfig struct {
	Contracts      []common.Address
	Topics         [][]common.Hash
	Confirmations  uint64
	Step           uint64
	ListingABI     abi.ABI
	AuctionABI     abi.ABI
	MarketplaceABI abi.ABI
}
