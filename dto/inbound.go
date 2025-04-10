package dto

import "encoding/json"

type Inbound struct {
	ID             int          `json:"id"`
	Up             int64        `json:"up"`
	Down           int64        `json:"down"`
	Total          int64        `json:"total"`
	Remark         string       `json:"remark"`
	Enable         bool         `json:"enable"`
	ExpiryTime     int64        `json:"expiryTime"`
	ClientStats    []ClientStat `json:"clientStats"`
	Listen         string       `json:"listen"`
	Port           int          `json:"port"`
	Protocol       string       `json:"protocol"`
	Settings       string       `json:"settings"`
	StreamSettings string       `json:"streamSettings"`
	Tag            string       `json:"tag"`
	Sniffing       string       `json:"sniffing"`
	Allocate       string       `json:"allocate"`
}

func (i Inbound) GetSettings() Settings {
	var settings Settings
	if err := json.Unmarshal([]byte(i.Settings), &settings); err != nil {
		panic(err)
	}
	return settings
}
