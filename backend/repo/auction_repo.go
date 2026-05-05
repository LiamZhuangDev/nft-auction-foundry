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

func (r *AuctionRepo) CreateAuction(a *models.Auction) error {
	return r.DB.Create(a).Error
}

func (r *AuctionRepo) GetAuctions() ([]models.Auction, error) {
	var auctions []models.Auction

	if err := r.DB.Find(&auctions).Error; err != nil {
		return nil, err
	}

	return auctions, nil
}

func (r *AuctionRepo) UpdateAuctionStatus(id uint64, finalPrice string, status bool) error {
	err := r.DB.
		Model(&models.Auction{}).
		Where("auction_id = ?", id).
		Updates(map[string]any{
			"final_price": finalPrice,
			"active":      status,
		}).Error

	return err
}
