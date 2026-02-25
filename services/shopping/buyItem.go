package shopping

type repository interface {
	BuyItem(id int, item string) error
}

type ShoppingService struct {
	repo repository
}

func NewShoppingService(repo repository) *ShoppingService {
	return &ShoppingService{
		repo: repo,
	}
}

func (sh *ShoppingService) Shop(id int, item string) error {
	return sh.repo.BuyItem(id, item)
}
