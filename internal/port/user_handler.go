package port

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
	"net/http"
)

func (h *HTTPServer) GetUserByID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var uriParam api.GetUserUriParameter
		err := ctx.ShouldBindUri(&uriParam)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		userID := uriParam.ID
		user, err := h.userService.GetUser(ctx, userID)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		accounts, err := h.userService.GetVisibleBalanceAccounts(ctx, userID)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		var accountsResp []api.GetBalanceAccountResponse
		for _, a := range accounts {
			accountResp := api.GetBalanceAccountResponse{
				Balance:  a.Balance,
				Currency: a.Currency,
				Type:     a.Type.String(),
				State:    a.State.String(),
				Metadata: unmarshallMetadata(a.Metadata),
			}
			accountsResp = append(accountsResp, accountResp)
		}

		userResp := api.GetUserResponse{
			ID:              user.ID,
			Email:           user.Email,
			Phone:           user.Phone,
			UserName:        user.ExternalID,
			Metadata:        unmarshallMetadata(user.Metadata),
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
			BalanceAccounts: accountsResp,
		}

		ctx.JSON(http.StatusOK, userResp)
	}
}

func (h *HTTPServer) AddUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req api.AddUserRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		user, err := h.userService.AddUser(ctx, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		resp := api.AddUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			UserName:  user.ExternalID,
			CreatedAt: user.CreatedAt,
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (h *HTTPServer) UpdateUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var uriParam api.UpdateUserUriParameter
		err := ctx.ShouldBindUri(&uriParam)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		var req api.UpdateUserRequest
		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		userID := uriParam.ID
		user, err := h.userService.UpdateUser(ctx, userID, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		resp := api.UpdateUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			UserName:  user.ExternalID,
			UpdatedAt: user.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// TODO this should be moved to utils
func unmarshallMetadata(bytes []byte) map[string]interface{} {
	var metaData map[string]interface{}
	json.Unmarshal(bytes, &metaData)
	return metaData
}
