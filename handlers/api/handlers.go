package api

import (
	"root/services/login"
	"root/services/sendCoins"
	"root/services/shopping"
	"root/services/userInfo"
)

type Handlers struct {
	loginService     *login.LoginService
	sendCoinsService *sendCoins.SendCoinsService
	shoppingService  *shopping.ShoppingService
	userinfoService  *userInfo.UserInfoService
}

func NewHandlers(lg *login.LoginService, sc *sendCoins.SendCoinsService, sh *shopping.ShoppingService, ui *userInfo.UserInfoService) *Handlers {
	return &Handlers{
		loginService:     lg,
		sendCoinsService: sc,
		shoppingService:  sh,
		userinfoService:  ui,
	}
}
