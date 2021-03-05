package domain

type ACTION string

const (
	TRY    ACTION = "TRY"
	COMMIT ACTION = "COMMIT"
	CANCEL ACTION = "CANCEL"
)

type Transaction struct {
	ID             uint64
	ParentID       uint64
	SrcAccountID   uint64
	DstAccountID   uint64
	SrcUserID      uint64
	DstUserID      uint64
	Amount         int64
	Currency       string
	SrcBalance     int64
	DstBalance     int64
	SrcAccountType string
	DstAccountType string
	Reversible     bool
	Metadata       TransactionMetadata
	CreatedAt      int64
}

type TransactionMetadata struct {
	Action          ACTION
	NextTransaction NextTransaction
}

type NextTransaction struct {
	SrcAccountID uint64
	SrcUserID    uint64
	DstAccountID uint64
	DstUserID    uint64
	Amount       int64
}
