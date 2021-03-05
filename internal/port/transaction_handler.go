package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
	"net/http"
)

// FixMe: Temporary for demo, delete this file later
func (h *HTTPServer) GetTransactionByUserID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var uriParam api.GetUserTransactionsUriParameter
		if err := ctx.ShouldBindUri(&uriParam); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		transactions, err := h.transactionService.GetTransactionByUserID(ctx, uriParam.UserID)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		transactionResponse := make([]api.TransactionResponse, len(transactions))
		for i, v := range transactions {
			transactionResponse[i] = api.TransactionResponse{
				ID:             v.ID,
				ParentID:       v.ParentID,
				SrcAccountID:   v.SrcAccountID,
				DstAccountID:   v.DstAccountID,
				Amount:         v.Amount,
				Currency:       v.Currency,
				SrcAccountType: v.SrcAccountType,
				DstAccountType: v.DstAccountType,
				CreatedAt:      v.CreatedAt,
			}
		}

		response := api.GetUserTransactionsResponse{
			Transactions: transactionResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
