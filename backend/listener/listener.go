package listener

import (
	"log"
	"nft-auction-backend/config"
	"nft-auction-backend/repo"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener(
	listingRepo *repo.ListingRepo,
	auctionRepo *repo.AuctionRepo,
	bidRepo *repo.BidRepo) error {
	log.Println("Listener Started")

	client, err := ethclient.Dial(config.EthNodeWsUrl)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Println("Subscribing to logs...")

	// Use a WaitGroup to keep the main goroutine alive while the listeners run in parallel
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := WatchListingEvents(client, listingRepo)
		if err != nil {
			log.Println("WatchListingEvents error: ", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := WatchAuctionEvents(client, auctionRepo, bidRepo)
		if err != nil {
			log.Println("WatchAuctionEvents error: ", err)
		}
	}()

	wg.Wait()

	log.Println("Listener Stopped")

	return nil
}
