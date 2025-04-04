package dto

type Settings struct {
	Clients    []Client      `json:"clients"`
	Decryption string        `json:"decryption"`
	Fallbacks  []interface{} `json:"fallbacks"`
}
