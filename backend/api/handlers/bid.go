package handlers

import (
	"math/big"
	"net/http"
	"nft-auction-backend/models"
	"nft-auction-backend/repo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBidsByAuction(repo *repo.BidRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		auctionId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction ID"})
			return
		}

		bids, err := repo.GetBidsByAuctionId(auctionId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": bids})
	}
}

func CreateBidForAuction(repo *repo.BidRepo) gin.HandlerFunc {
	type CreateBidRequest struct {
		Bidder string `json:"bidder" binding:"required"`
		Amount string `json:"amount" binding:"required"`
	}
	return func(c *gin.Context) {
		// 1. Parse auction ID from URL
		idStr := c.Param("id")
		auctionId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction ID"})
			return
		}

		// 2. Parse JSON body
		var req CreateBidRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// 3. Basic validation
		amount := new(big.Int)
		_, ok := amount.SetString(req.Amount, 10)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
			return
		}

		if amount.Sign() <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
			return
		}

		// 4. Create bid
		bid := &models.Bid{
			AuctionID: auctionId,
			Bidder:    req.Bidder,
			Amount:    req.Amount,
			Timestamp: uint64(time.Now().Unix()),
		}
		err = repo.CreateBid(bid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create bid"})
			return
		}

		c.JSON(201, gin.H{"data": bid})
	}
}
