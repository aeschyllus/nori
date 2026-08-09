package account

type AccountResponse struct {
	ID     int64  `json:"accountId"`
	Name   string `json:"name"`
	Amount int64  `json:"amountInCents"`
}

type AccountRequest struct {
	Name   string `json:"name" binding:"required"`
	Amount int64  `json:"amountInCents"`
}
