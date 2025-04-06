package dto

import "encoding/json"

type Client struct {
	Comment    string      `json:"comment"`
	Email      string      `json:"email"`
	Enable     bool        `json:"enable"`
	ExpiryTime int64       `json:"expiryTime"`
	Flow       string      `json:"flow"`
	Id         string      `json:"id"`
	LimitIp    int         `json:"limitIp"`
	Reset      int         `json:"reset"`
	SubId      string      `json:"subId"`
	TgId       json.Number `json:"tgId"`
	TotalGB    int64       `json:"totalGB"`
}
