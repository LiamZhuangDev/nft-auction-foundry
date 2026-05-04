package repo

import (
	"nft-auction-backend/models"

	"gorm.io/gorm"
)

type ListingRepo struct {
	DB *gorm.DB
}

func NewListingRepo(db *gorm.DB) *ListingRepo {
	return &ListingRepo{DB: db}
}

func (r *ListingRepo) CreateListing(l *models.Listing) error {
	return r.DB.Create(l).Error
}

func (r *ListingRepo) GetListings() ([]models.Listing, error) {
	var listings []models.Listing

	if err := r.DB.Find(&listings).Error; err != nil {
		return nil, err
	}

	return listings, nil
}
