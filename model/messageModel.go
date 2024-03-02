package model

type Message struct {
	Pengguna Pengguna `json:"Pengguna"`
	Message  string   `json:"Message"`
}
