package models

type UserActivity struct {
	Coins       int          `json:"coins"`
	Inventory   []*Inventory `json:"inventory"`
	CoinHistory *CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Quantity int    `json:"quantity"`
	Type     string `json:"type"`
}

type CoinHistory struct {
	Received []*Received `json:"received"`
	Sent     []*Sent     `json:"sent"`
}

type Received struct {
	Amount   int    `json:"amount"`
	FromUser string `json:"fromUser"`
}

type Sent struct {
	Amount int    `json:"amount"`
	ToUser string `json:"toUser"`
}
