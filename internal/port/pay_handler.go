package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
	"net/http"
)

func (h *HTTPServer) TryPay() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var tryPayRequest api.TryPayRequest
		var idemHeader api.IdemHeader

		if err := ctx.ShouldBindJSON(&tryPayRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		if err := ctx.ShouldBindHeader(&idemHeader); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		parentId, err := h.payService.TryPay(ctx, idemHeader.IdemKey, tryPayRequest)

		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		tryPayResp := api.TryPayResponse{
			ID: parentId,
		}

		h.RespondWithOK(ctx, tryPayResp)
	}
}

func (h *HTTPServer) CommitPay() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var commitPayRequest api.CommitPayRequest
		if err := ctx.BindJSON(&commitPayRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		if err := h.payService.CommitPay(ctx, commitPayRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

func (h *HTTPServer) CancelPay() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var cancelPayRequest api.CancelPayRequest
		if err := ctx.BindJSON(&cancelPayRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		if err := h.payService.CancelPay(ctx, cancelPayRequest); err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}
