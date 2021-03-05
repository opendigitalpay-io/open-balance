package api

// FixMe: Temporary for demo, delete this file later

type GetUserTransactionsUriParameter struct {
	UserID uint64 `uri:"userId" binding:"required"`
}

type GetUserTransactionsResponse struct {
	Transactions []TransactionResponse
}

type TransactionResponse struct {
	ID             uint64
	ParentID       uint64
	SrcAccountID   uint64
	DstAccountID   uint64
	Amount         int64
	Currency       string
	SrcAccountType string
	DstAccountType string
	CreatedAt      int64
}
