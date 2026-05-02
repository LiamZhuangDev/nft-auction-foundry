package api

import (
	"nft-auction-backend/api/handlers"
	"nft-auction-backend/repo"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	listingRepo *repo.ListingRepo,
	auctionRepo *repo.AuctionRepo,
	bidRepo *repo.BidRepo,
) {
	r.GET("/listings", handlers.GetListings(listingRepo))
	r.GET("/auctions", handlers.GetAuctions(auctionRepo))
	r.GET("/bids/:auction_id", handlers.GetBids(bidRepo))
}
