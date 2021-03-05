package api

type TryPayRequest struct {
	UserID       uint64                 `json:"userId" binding:"required"`
	BusinessID   uint64                 `json:"businessId" binding:"required"`
	Amount       int64                  `json:"amount" binding:"required"`
	ChargeAmount int64                  `json:"chargeAmount" binding:"gte=0"`
	Currency     string                 `json:"currency" binding:"required"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type TryPayResponse struct {
	ID uint64 `json:"id"`
}

type CommitPayRequest struct {
	ParentID uint64 `json:"parentId" binding:"required"`
	Metadata map[string]interface{}
}

type CancelPayRequest struct {
	ParentID uint64 `json:"parentId" binding:"required"`
	Metadata map[string]interface{}
}
