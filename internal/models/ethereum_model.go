package models

type SendEthRequest struct {
	ToAddress string  `json:"toAddress"`
	Amount    float64 `json:"amount"`
}
