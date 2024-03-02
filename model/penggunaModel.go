package model

type Pengguna struct {
	Id          int    `json:"idWallet"`
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	Password    string `json:"Password"`
	Type        string `json:"Type"`
	DisableUser bool   `json:"DisableUser"`
	Activated   bool   `json:"Activated"`
}

type PenggunaResponse struct {
	Status   int      `json:"status"`
	Message  string   `json:"message"`
	Wallet   Wallet   `json:"Wallet"`
	Pengguna Pengguna `json:"Pengguna"`
}

type OTPModel struct {
	Pengguna Pengguna `json:"Pengguna"`
	OTP      string   `json:"OTP"`
}
