package models

type Listing struct {
	ID          uint64 `gorm:"primaryKey"`  // DB internal ID
	ListingID   uint64 `gorm:"uniqueIndex"` // on-chain ID
	Seller      string
	NftContract string
	TokenId     uint64
	Active      bool
}
