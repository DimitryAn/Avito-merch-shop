package userInfo

import (
	"root/models"
)

type repopository interface {
	GetAllActivity(id int) (*models.UserActivity, error)
}

type UserInfoService struct {
	repo repopository
}

func NewUsernInfoService(repo repopository) *UserInfoService {
	return &UserInfoService{
		repo: repo,
	}
}

func (ui *UserInfoService) Activity(id int) (*models.UserActivity, error) {
	return ui.repo.GetAllActivity(id)
}
