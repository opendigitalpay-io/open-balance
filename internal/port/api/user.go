package api

type GetUserUriParameter struct {
	ID uint64 `uri:"id" binding:"required"`
}

type GetUserResponse struct {
	ID              uint64                      `json:"id"`
	Email           string                      `json:"email"`
	Phone           string                      `json:"phone"`
	UserName        string                      `json:"userName"`
	Metadata        map[string]interface{}      `json:"metadata"`
	CreatedAt       int64                       `json:"createdAt"`
	UpdatedAt       int64                       `json:"updatedAt"`
	BalanceAccounts []GetBalanceAccountResponse `json:"balanceAccounts"`
}

type AddUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	UserName string `json:"userName" binding:"required"`
}

type AddUserResponse struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserName  string `json:"userName"`
	CreatedAt int64  `json:"createdAt"`
}

type UpdateUserUriParameter struct {
	ID uint64 `json:"id" binding:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	UserName string `json:"userName"`
}

type UpdateUserResponse struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserName  string `json:"userName"`
	UpdatedAt int64  `json:"updatedAt"`
}
