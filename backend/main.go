package main

import (
	"log"
	"nft-auction-backend/api"
	"nft-auction-backend/config"
	"nft-auction-backend/listener"
	"nft-auction-backend/repo"
)

func main() {
	log.Println("Starting app...")

	db := config.ConnectDB()

	listingRepo := repo.NewListingRepo(db)
	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)

	log.Println("Database connected")

	// Start ethereum listener (background)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Listener crashed: ", r)
			}
		}()
		listener.StartListener(listingRepo, auctionRepo, bidRepo)
	}()

	// Start http server (blocking)
	api.StartServer(listingRepo, auctionRepo, bidRepo)
}
