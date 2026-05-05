package listener

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"nft-auction-backend/config"
	"nft-auction-backend/models"
	"nft-auction-backend/repo"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WatchAuctionEvents(client *ethclient.Client, auctionRepo *repo.AuctionRepo, bidRepo *repo.BidRepo) error {
	// Load the contract ABI
	file, err := os.ReadFile("abi/auction.json")
	if err != nil {
		return err
	}

	var abiJson struct {
		ABI json.RawMessage `json:"abi"`
	}

	err = json.Unmarshal(file, &abiJson)
	if err != nil {
		return err
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJson.ABI)))
	if err != nil {
		return err
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{config.AuctionAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	log.Println("Listening for Auction Events...")

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			handleAuctionLog(abi, vLog, auctionRepo, bidRepo)
		}
	}
}

func handleAuctionLog(abi abi.ABI, vLog types.Log, auctionRepo *repo.AuctionRepo, bidRepo *repo.BidRepo) {
	event, err := abi.EventByID(vLog.Topics[0])
	if err != nil {
		log.Printf("Failed to extract event: %v\n", err)
		return
	}

	switch event.Name {
	case "AuctionCreated":
		log.Println("Received AuctionCreated event")

		// Unpack the non-indexed event data into the struct
		var data struct {
			TokenId    *big.Int
			StartPrice *big.Int
			EndTime    *big.Int
		}

		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v\n", err)
			return
		}

		// Extract indexed params from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		seller := common.BytesToAddress(vLog.Topics[2].Bytes()) // extract the last 20 bytes for address
		nftContract := common.BytesToAddress(vLog.Topics[3].Bytes())

		// Save auction
		err = auctionRepo.CreateAuction(&models.Auction{
			AuctionID:   auctionId.Uint64(),
			Seller:      seller.Hex(),
			NftContract: nftContract.Hex(),
			TokenId:     data.TokenId.String(),
			StartPrice:  data.StartPrice.String(),
			EndTime:     data.EndTime.Uint64(),
			Active:      true,
		})
		if err != nil {
			log.Printf("Failed to create auction: %v\n", err)
			return
		}

		log.Printf("AuctionCreated: AuctionID=%d, Seller=%s, NFT=%s, TokenID=%s, StartPrice=%s, EndTime=%d\n",
			auctionId.Uint64(), seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.StartPrice.String(), data.EndTime.Uint64())
	case "BidPlaced":
		log.Println("Received BidPlaced event")

		// Unpack the non-indexed event data into the struct
		var data struct {
			Amount *big.Int
		}

		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack data: %v\n", err)
			return
		}

		// Extract indexed params from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		bidder := common.BytesToAddress(vLog.Topics[2].Bytes())

		// Save bid
		err = bidRepo.CreateBid(&models.Bid{
			AuctionID: auctionId.Uint64(),
			Bidder:    bidder.Hex(),
			Amount:    data.Amount.String(),
			Timestamp: uint64(time.Now().Unix()),
		})
		if err != nil {
			log.Printf("Failed to create bid, err: %v\n", err)
			return
		}

		log.Printf("BidPlaced: AuctionId=%d, Bidder=%s, Amount=%s\n", auctionId.Uint64(), bidder.Hex(), data.Amount.String())
	case "AuctionEnded":
		log.Println("Receive AuctionEnded event")

		// Unpack non-index event data into the struct
		var data struct {
			Amount *big.Int
		}
		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v\n", err)
			return
		}

		// Extract indexed params from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		winner := common.BytesToAddress(vLog.Topics[2].Bytes())

		// Update Auction
		err = auctionRepo.UpdateAuctionStatus(auctionId.Uint64(), false)
		if err != nil {
			log.Printf("Failed to update auction %d status, error: %v\n", auctionId.Uint64(), err)
			return
		}

		log.Printf("AuctionEnded: AuctionId=%d, Winner=%s, Amount=%s\n", auctionId.Uint64(), winner.Hex(), data.Amount.String())
	default:
		log.Printf("Unknown event: %s\n", vLog.Topics[0].Hex())
	}
}
