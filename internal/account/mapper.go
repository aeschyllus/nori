package account

func toAccountResponse(a Account) AccountResponse {
	return AccountResponse{
		ID:     a.ID,
		Name:   a.Name,
		Amount: a.Amount,
	}
}

func toAccountResponses(accounts []Account) []AccountResponse {
	resp := make([]AccountResponse, len(accounts))
	for i, a := range accounts {
		resp[i] = toAccountResponse(a)
	}
	return resp
}

func toAccountFromRequest(req AccountRequest) Account {
	return Account{
		Name:   req.Name,
		Amount: req.Amount,
	}
}
