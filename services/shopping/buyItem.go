package shopping

import "context"

type repository interface {
	BuyItem(ctx context.Context, id int, item string) error
}

type ShoppingService struct {
	repo repository
}

func NewShoppingService(repo repository) *ShoppingService {
	return &ShoppingService{
		repo: repo,
	}
}

func (sh *ShoppingService) Shop(ctx context.Context, id int, item string) error {
	return sh.repo.BuyItem(ctx, id, item)
}
