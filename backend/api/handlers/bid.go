package handlers

import (
	"net/http"
	"nft-auction-backend/repo"

	"github.com/gin-gonic/gin"
)

func GetBids(repo *repo.BidRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		listings, err := repo.GetBids()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": listings})
	}
}
