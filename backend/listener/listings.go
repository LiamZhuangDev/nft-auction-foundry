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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WatchListingEvents(client *ethclient.Client, listingRepo *repo.ListingRepo) error {
	file, err := os.ReadFile("abi/marketplace.json")
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
		Addresses: []common.Address{config.MarketplaceAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	log.Println("Listening for Listing events...")

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			handleListingLog(abi, vLog, listingRepo)
		}
	}
}

func handleListingLog(abi abi.ABI, vLog types.Log, listingRepo *repo.ListingRepo) {
	event, err := abi.EventByID(vLog.Topics[0])
	if err != nil {
		log.Printf("Unknown event: %s\n", vLog.Topics[0].Hex())
		return
	}

	switch event.Name {
	case "Listed":
		log.Println("Handling Listed event...")

		// Unpack the non-index event data into the struct
		var data struct {
			TokenId   *big.Int
			ListingId *big.Int
		}

		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v\n", err)
			return
		}

		// Extract indexed params from the topics
		seller := common.HexToAddress(vLog.Topics[1].Hex())
		nftContract := common.HexToAddress(vLog.Topics[2].Hex())

		err = listingRepo.CreateListing(&models.Listing{
			Seller:      seller.Hex(),
			NftContract: nftContract.Hex(),
			TokenId:     data.TokenId.String(),
			Active:      true,
		})

		if err != nil {
			log.Printf("Failed to list NFT: %v\n", err)
			return
		}

		log.Printf("New listing: Seller=%s, NFT contract=%s, Token ID=%s, Listing ID=%s",
			seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.ListingId.String())
	default:
		log.Printf("Unhandled event: %s", event.Name)
	}
}
