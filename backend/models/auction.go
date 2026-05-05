package models

type Auction struct {
	ID          uint64 `gorm:"primaryKey"`  // DB internal ID
	AuctionID   uint64 `gorm:"uniqueIndex"` // on-chain ID
	Seller      string
	NftContract string
	TokenId     uint64
	StartPrice  string
	FinalPrice  string
	EndTime     uint64
	Active      bool
	// By default, GORM would assume bids.auction_id -> auctions.id
	// But it's linking bids using on-chain ID. So when access Auction.Bids,
	// use `references:AuctionID` to match bids.auction_id with auctions.auction_id
	Bids []Bid `gorm:"foreignKey:AuctionID;references:AuctionID"`
}
