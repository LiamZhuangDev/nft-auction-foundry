package main

import (
	"context"
	"encoding/json"
	"log"
	"nft-auction-backend/api"
	"nft-auction-backend/config"
	"nft-auction-backend/repo"
	"nft-auction-backend/service"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	log.Println("Starting app...")

	db := config.ConnectDB()

	listingRepo := repo.NewListingRepo(db)
	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)
	eventRepo := repo.NewEventRepo(db)

	log.Println("Database connected")

	// Start ethereum listener (background)
	// go func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			log.Println("Listener crashed: ", r)
	// 		}
	// 	}()
	// 	listener.StartListener(listingRepo, auctionRepo, bidRepo)
	// }()

	// Start ethereum poller
	go func() {
		client, err := ethclient.Dial(config.EthNodeWsUrl)
		if err != nil {
			log.Println("Failed to connect to ethereum")
		}
		defer client.Close()

		auctionABI, err := initABI("abi/auction.json")
		if err != nil {
			panic(err)
		}

		marketpalceABI, err := initABI("abi/marketplace.json")
		if err != nil {
			panic(err)
		}

		// Topic[0] = Listed OR AuctionCreated OR BidPlaced OR AuctionEnded
		// It means poller will be configured to monitor these events
		topics := [][]common.Hash{
			{
				marketpalceABI.Events["Listed"].ID,
				auctionABI.Events["AuctionCreated"].ID,
				auctionABI.Events["BidPlaced"].ID,
				auctionABI.Events["AuctionEnded"].ID,
			},
		}

		cfg := config.PollerConfig{
			Contracts: []common.Address{
				config.AuctionAddress,
				config.MarketplaceAddress,
			},
			Topics:         topics,
			Confirmations:  6,
			Step:           2000,
			AuctionABI:     auctionABI,
			MarketplaceABI: marketpalceABI,
		}
		poller := service.NewPoller(client, cfg, eventRepo, listingRepo, auctionRepo, bidRepo)
		poller.Start(context.Background())
	}()

	// Start http server (blocking)
	api.StartServer(listingRepo, auctionRepo, bidRepo)
}

func initABI(filePath string) (abi.ABI, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return abi.ABI{}, err
	}

	// artifact structure
	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}

	if err := json.Unmarshal(file, &artifact); err != nil {
		return abi.ABI{}, err
	}

	return abi.JSON(strings.NewReader(string(artifact.ABI)))
}
