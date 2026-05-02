package config

import (
	"log"
	"nft-auction-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// See readme for database setup
func ConnectDB() *gorm.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/nft_marketplace?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Listing{}, &models.Auction{}, &models.Bid{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
