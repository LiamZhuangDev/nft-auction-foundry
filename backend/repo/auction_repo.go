package repo

type AuctionRepo struct {
}

func (*AuctionRepo) GetAuctions() ([]any, error) {
	return []any{}, nil
}
