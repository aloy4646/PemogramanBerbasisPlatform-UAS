package model

type Transaction struct {
	Id          int    `json:"id"`
	IdWallet    string `json:"idWallet"`
	Datetime    string `json:"datetime"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

type DetailTransaction struct {
	Id           int           `json:"id"`
	Currency     string        `json:"currency"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	DisableUser  bool          `json:"disableUser"`
	Transactions []Transaction `json:"transaction"`
}
