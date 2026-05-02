package handlers

import (
	"net/http"
	"nft-auction-backend/repo"

	"github.com/gin-gonic/gin"
)

func GetAuctions(repo *repo.AuctionRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		listings, err := repo.GetAuctions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": listings})
	}
}
