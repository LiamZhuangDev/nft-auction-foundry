package service

import (
	"context"
	"log"
	"math/big"
	"nft-auction-backend/config"
	"nft-auction-backend/models"
	"nft-auction-backend/repo"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Poller struct {
	Client         *ethclient.Client
	Contracts      []common.Address
	Topics         [][]common.Hash
	Confirmations  uint64
	Step           uint64
	eventRepo      *repo.EventRepo
	listingRepo    *repo.ListingRepo
	auctionRepo    *repo.AuctionRepo
	bidRepo        *repo.BidRepo
	auctionABI     abi.ABI
	marketplaceABI abi.ABI
}

func NewPoller(
	client *ethclient.Client,
	cfg config.PollerConfig,
	eventRepo *repo.EventRepo,
	listingRepo *repo.ListingRepo,
	auctionRepo *repo.AuctionRepo,
	bidRepo *repo.BidRepo,
) *Poller {
	return &Poller{
		Client:         client,
		Contracts:      cfg.Contracts,
		Topics:         cfg.Topics,
		Confirmations:  cfg.Confirmations,
		Step:           cfg.Step,
		eventRepo:      eventRepo,
		listingRepo:    listingRepo,
		auctionRepo:    auctionRepo,
		bidRepo:        bidRepo,
		auctionABI:     cfg.AuctionABI,
		marketplaceABI: cfg.MarketplaceABI,
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.sync(ctx); err != nil {
				log.Printf("Sync error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Poller) sync(ctx context.Context) error {
	latest, err := p.Client.BlockNumber(ctx)
	log.Println("Latest block number: ", latest)

	if err != nil {
		return err
	}

	if latest < p.Confirmations {
		return nil
	}

	safeBlock := latest - p.Confirmations
	lastProcessed, err := p.eventRepo.GetLastProcessedBlock(ctx)
	if err != nil {
		return err
	}

	if safeBlock <= lastProcessed {
		return nil
	}

	from := lastProcessed + 1
	to := safeBlock

	log.Printf("Syncing blocks [%d -> %d]", from, to)

	return p.processRange(ctx, from, to)
}

func (p *Poller) processRange(ctx context.Context, from, to uint64) error {
	for start := from; start <= to; start += p.Step {
		end := min(start+p.Step-1, to)

		if err := p.processBatch(ctx, start, end); err != nil {
			return err
		}

		if err := p.eventRepo.SetLastProcessedBlock(ctx, end); err != nil {
			return err
		}
	}

	return nil
}

func (p *Poller) processBatch(ctx context.Context, from, to uint64) error {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: p.Contracts,
		Topics:    p.Topics,
	}

	logs, err := p.Client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, log := range logs {
		if err := p.handleLog(log); err != nil {
			return err
		}
	}

	return nil
}

func (p *Poller) handleLog(vLog types.Log) error {
	switch vLog.Topics[0] {
	case p.marketplaceABI.Events["Listed"].ID:
		return p.handleNFTListed(vLog)
	case p.auctionABI.Events["AuctionCreated"].ID:
		return p.handleAuctionCreated(vLog)
	case p.auctionABI.Events["BidPlaced"].ID:
		return p.handleBidPlaced(vLog)
	case p.auctionABI.Events["AuctionEnded"].ID:
		return p.handleAuctionFinalize(vLog)
	default:
		log.Printf("Unknown event: %s\n", vLog.Topics[0].Hex())
	}

	return nil
}

func (p *Poller) handleNFTListed(vLog types.Log) error {
	log.Println("NFT Listed Event")

	// Unpack the non-index event data into the struct
	var data struct {
		TokenId   *big.Int
		ListingId *big.Int
	}

	err := p.marketplaceABI.UnpackIntoInterface(&data, "Listed", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log: %v\n", err)
		return err
	}

	// Extract indexed params from the topics
	seller := common.HexToAddress(vLog.Topics[1].Hex())
	nftContract := common.HexToAddress(vLog.Topics[2].Hex())

	err = p.listingRepo.CreateListing(&models.Listing{
		ListingID:   data.ListingId.Uint64(),
		Seller:      seller.Hex(),
		NftContract: nftContract.Hex(),
		TokenId:     data.TokenId.Uint64(),
		Active:      true,
	})

	if err != nil {
		log.Printf("Failed to list NFT: %v\n", err)
		return err
	}

	log.Printf("New listing: Seller=%s, NFT contract=%s, Token ID=%s, Listing ID=%s",
		seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.ListingId.String())

	return nil
}

func (p *Poller) handleAuctionCreated(vLog types.Log) error {
	log.Println("Auction Created Event")

	// Unpack the non-indexed event data into the struct
	var data struct {
		TokenId    *big.Int
		StartPrice *big.Int
		End        *big.Int
	}

	err := p.auctionABI.UnpackIntoInterface(&data, "AuctionCreated", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log: %v\n", err)
		return err
	}

	// Extract indexed params from the topics
	auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
	seller := common.BytesToAddress(vLog.Topics[2].Bytes()) // extract the last 20 bytes for address
	nftContract := common.BytesToAddress(vLog.Topics[3].Bytes())

	// Save auction
	err = p.auctionRepo.CreateAuction(&models.Auction{
		AuctionID:   auctionId.Uint64(),
		Seller:      seller.Hex(),
		NftContract: nftContract.Hex(),
		TokenId:     data.TokenId.Uint64(),
		StartPrice:  data.StartPrice.String(),
		EndTime:     data.End.Uint64(),
		Active:      true,
	})
	if err != nil {
		log.Printf("Failed to create auction: %v\n", err)
		return err
	}

	log.Printf("AuctionCreated: AuctionID=%d, Seller=%s, NFT=%s, TokenID=%s, StartPrice=%s, EndTime=%d\n",
		auctionId.Uint64(), seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.StartPrice.String(), data.End.Uint64())

	return nil
}

func (p *Poller) handleBidPlaced(vLog types.Log) error {
	log.Println("Bid Placed Event")

	// Unpack the non-indexed event data into the struct
	var data struct {
		Amount *big.Int
	}

	err := p.auctionABI.UnpackIntoInterface(&data, "BidPlaced", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack data: %v\n", err)
		return err
	}

	// Extract indexed params from the topics
	auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
	bidder := common.BytesToAddress(vLog.Topics[2].Bytes())

	// Save bid
	err = p.bidRepo.CreateBid(&models.Bid{
		AuctionID: auctionId.Uint64(),
		Bidder:    bidder.Hex(),
		Amount:    data.Amount.String(),
		Timestamp: uint64(time.Now().Unix()),
	})
	if err != nil {
		log.Printf("Failed to create bid, err: %v\n", err)
		return err
	}

	log.Printf("BidPlaced: AuctionId=%d, Bidder=%s, Amount=%s\n", auctionId.Uint64(), bidder.Hex(), data.Amount.String())

	return nil
}

func (p *Poller) handleAuctionFinalize(vLog types.Log) error {
	log.Println("Auction Ended Event")

	// Unpack non-index event data into the struct
	var data struct {
		Amount *big.Int
	}
	err := p.auctionABI.UnpackIntoInterface(&data, "AuctionEnded", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack log: %v\n", err)
		return err
	}

	// Extract indexed params from the topics
	auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
	winner := common.BytesToAddress(vLog.Topics[2].Bytes())

	// Update Auction
	err = p.auctionRepo.UpdateAuctionStatus(auctionId.Uint64(), data.Amount.String(), false)
	if err != nil {
		log.Printf("Failed to update auction %d status, error: %v\n", auctionId.Uint64(), err)
		return err
	}

	log.Printf("AuctionEnded: AuctionId=%d, Winner=%s, Amount=%s\n", auctionId.Uint64(), winner.Hex(), data.Amount.String())

	return nil
}
