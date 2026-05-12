package sendCoins

import (
	"context"
	"root/models"
)

type repository interface {
	SendCoins(ctx context.Context, fromId int, toId int, amount int) error
	GetUserByName(ctx context.Context, u *models.User) error
}

type SendCoinsService struct {
	repo repository
}

func NewSendCoinsService(repo repository) *SendCoinsService {
	return &SendCoinsService{
		repo: repo,
	}
}

func (scs *SendCoinsService) SendCoins(ctx context.Context, fromId int, toUser string, amount int) error {
	toU := models.User{
		Username: toUser,
	}

	err := scs.repo.GetUserByName(ctx, &toU) // получить id кому отправляем
	if err != nil {
		return err
	}

	return scs.repo.SendCoins(ctx, fromId, toU.ID, amount)
}
