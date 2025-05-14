package models

type SendEthRequest struct {
	ToAddress string  `json:"to_Address"`
	Amount    float64 `json:"amount"`
}

type GetEthRequest struct {
	UserId string `json:"user_id"`
}

type EthBalanceResponse struct {
	BalanceETH string `json:"balance_eth"`
	BalanceUSD string `json:"balance_usd"`
}
