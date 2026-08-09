package account

// Account is a single money account.
// Amount is the balance in integer cents; a negative value represents an overdraft.
type Account struct {
	ID     int64
	Name   string
	Amount int64
}
