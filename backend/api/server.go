package api

import (
	"nft-auction-backend/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer(listingRepo *repo.ListingRepo, auctionRepo *repo.AuctionRepo, bidRepo *repo.BidRepo) {
	r := gin.Default()
	r.Use(cors.Default())

	RegisterRoutes(r, listingRepo, auctionRepo, bidRepo)

	r.Run(":8081")
}
