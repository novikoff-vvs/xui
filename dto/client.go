package dto

type Client []struct {
	Comment    string `json:"comment"`
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	ExpiryTime int    `json:"expiryTime"`
	Flow       string `json:"flow"`
	Id         string `json:"id"`
	LimitIp    int    `json:"limitIp"`
	Reset      int    `json:"reset"`
	SubId      string `json:"subId"`
	TgId       int    `json:"tgId"`
	TotalGB    int    `json:"totalGB"`
}
