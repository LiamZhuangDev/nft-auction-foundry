package repo

type ListingRepo struct {
}

func (*ListingRepo) GetAllListings() ([]any, error) {
	return []any{}, nil
}
