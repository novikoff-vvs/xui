package requests

import "github.com/novikoff-vvs/xui/dto"

type AddClientToInboundRequest struct {
	InboundId int    `json:"id"`
	Settings  string `json:"settings"`
}

type AddClientToInboundClientRequest struct {
	Clients []dto.Client `json:"clients"`
}
