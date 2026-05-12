package userInfo

import (
	"context"
	"root/models"
)

type repopository interface {
	GetAllActivity(ctx context.Context, id int) (*models.UserActivity, error)
}

type UserInfoService struct {
	repo repopository
}

func NewUsernInfoService(repo repopository) *UserInfoService {
	return &UserInfoService{
		repo: repo,
	}
}

func (ui *UserInfoService) Activity(ctx context.Context, id int) (*models.UserActivity, error) {
	return ui.repo.GetAllActivity(ctx, id)
}
