package idem

type Idempotency struct {
	ID          uint64
	IdemID      string
	IsCompleted bool
	Response    string
	CreatedAt   int64
	UpdatedAt   int64
}
