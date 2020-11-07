package models

import "time"

// TwitchDatas represents the data content of twitch API response
type TwitchDatas struct {
	Title     string    `json:"title"`
	StartedAt time.Time `json:"started_at"`
}

// TwitchResponse represents the json response from twitch API
type TwitchResponse struct {
	Datas []TwitchDatas `json:"data"`
}
