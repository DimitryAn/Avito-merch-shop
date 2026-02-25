package sendCoins

import (
	"root/models"
)

type repository interface {
	SendCoins(fromId int, toId int, amount int) error
	GetUserByName(u *models.User) error
}

type SendCoinsService struct {
	repo repository
}

func NewSendCoinsService(repo repository) *SendCoinsService {
	return &SendCoinsService{
		repo: repo,
	}
}

func (scs *SendCoinsService) SendCoins(fromId int, toUser string, amount int) error {
	toU := models.User{
		Username: toUser,
	}

	err := scs.repo.GetUserByName(&toU) // получить id кому отправляем
	if err != nil {
		return err
	}

	return scs.repo.SendCoins(fromId, toU.ID, amount)
}
