package models

type Bid struct {
	ID        uint64 `gorm:"primaryKey"`
	AuctionID uint64 `gorm:"index"`
	Bidder    string
	Amount    string
	Timestamp uint64
}
