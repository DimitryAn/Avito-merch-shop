package errs

import "errors"

var (
	WrongPassword  = errors.New("wrong password")
	NotEnoughMoney = errors.New("not enough money")
)
