package api

type TryTopUpRequest struct {
	UserID   uint64                 `json:"userId" binding:"required"`
	Amount   int64                  `json:"amount" binding:"required"`
	Currency string                 `json:"currency" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

type TryTopUpResponse struct {
	ID uint64 `json:"id"`
}

type CommitTopUpRequest struct {
	ParentID uint64                 `json:"parentId" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

type CancelTopUpRequest struct {
	ParentID uint64                 `json:"parentId" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}
