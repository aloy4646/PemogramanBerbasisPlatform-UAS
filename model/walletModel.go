package model

type Wallet struct {
	Id          int    `json:"id"`
	Currency    string `json:"currency"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisableUser bool   `json:"disableUser"`
}
