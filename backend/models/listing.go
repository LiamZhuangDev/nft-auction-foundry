package models

type Listing struct {
	ID          uint64 `gorm:"primaryKey"`
	Seller      string
	NftContract string
	TokenId     string
	Active      bool
}
