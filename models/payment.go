package models

type Payment struct {
	Amount    float64 `json:"amount"`
	ID        string  `json:"correlationId"`
	RequestAt string  `json:"requested_at"`
}
