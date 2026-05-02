package repo

import (
	"nft-auction-backend/models"

	"gorm.io/gorm"
)

type BidRepo struct {
	DB *gorm.DB
}

func NewBidRepo(db *gorm.DB) *BidRepo {
	return &BidRepo{DB: db}
}

func (r *BidRepo) GetBids(auctionID uint64) ([]models.Bid, error) {
	var bids []models.Bid

	err := r.DB.
		Where("auction_id = ?", auctionID).
		Order("amount DESC").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}
