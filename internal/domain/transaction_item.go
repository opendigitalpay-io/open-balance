package domain

type TransactionItem struct {
	SrcAccount BalanceAccount
	SrcUserID  uint64
	DstAccount BalanceAccount
	DstUserID  uint64
	Amount     int64
}
