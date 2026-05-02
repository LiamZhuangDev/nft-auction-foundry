package main

import (
	"log"
	"nft-auction-backend/api"
	"nft-auction-backend/config"
	"nft-auction-backend/repo"
)

func main() {
	log.Println("Starting app...")

	db := config.ConnectDB()

	listingRepo := repo.NewListingRepo(db)
	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)

	log.Println("Database connected")

	api.StartServer(listingRepo, auctionRepo, bidRepo)
}
