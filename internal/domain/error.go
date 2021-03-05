package domain

type TransactionError struct {
	What string
}

func (e TransactionError) Error() string {
	return e.What
}

type IdemError struct {
	What string
}

func (e IdemError) Error() string {
	return e.What
}
