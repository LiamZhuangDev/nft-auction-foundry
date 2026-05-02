package repo

import (
	"nft-auction-backend/models"

	"gorm.io/gorm"
)

type AuctionRepo struct {
	DB *gorm.DB
}

func NewAuctionRepo(db *gorm.DB) *AuctionRepo {
	return &AuctionRepo{DB: db}
}

func (r *AuctionRepo) GetAuctions() ([]models.Auction, error) {
	var auctions []models.Auction

	if err := r.DB.Find(&auctions).Error; err != nil {
		return nil, err
	}

	return auctions, nil
}
