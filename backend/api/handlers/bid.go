package handlers

import (
	"net/http"
	"nft-auction-backend/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBids(repo *repo.BidRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		auctionId := c.Param("auction_id")
		id, err := strconv.ParseUint(auctionId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction ID"})
			return
		}

		listings, err := repo.GetBids(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": listings})
	}
}
