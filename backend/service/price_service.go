package service

import (
	"math/big"
	"nft-auction-backend/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

type PriceService struct {
	feed *Chainlink
}

func NewPriceService(client *ethclient.Client) (*PriceService, error) {
	feed, err := NewChainlink(config.ChainlinkFeedAddress, client)
	if err != nil {
		return nil, err
	}

	return &PriceService{
		feed: feed,
	}, nil
}

func (p *PriceService) GetETHUSD() (float64, error) {
	data, err := p.feed.LatestRoundData(nil)
	if err != nil {
		return 0, err
	}

	price := new(big.Float).Quo(
		new(big.Float).SetInt(data.Answer),
		big.NewFloat(1e8), // Chainlink ETH/USD decimals
	)

	value, _ := price.Float64()

	return value, nil
}

func (p *PriceService) WeiToUSD(wei *big.Int) (float64, error) {
	ethUsd, err := p.GetETHUSD()
	if err != nil {
		return 0, err
	}

	ethVal := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18), // ETH decimals
	)

	ethFloat, _ := ethVal.Float64()

	return ethFloat * ethUsd, nil
}
