package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
)

func (h *HTTPServer) TryTopUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var header api.IdemHeader
		err := ctx.ShouldBindHeader(&header)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		var req api.TryTopUpRequest
		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		parentID, err := h.topUpService.TryTopUp(ctx, header.IdemKey, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		resp := api.TryTopUpResponse{
			ID: parentID,
		}

		h.RespondWithOK(ctx, resp)
	}
}

func (h *HTTPServer) CommitTopUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req api.CommitTopUpRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		err = h.topUpService.CommitTopUp(ctx, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		h.RespondWithOK(ctx, nil)
	}
}

func (h *HTTPServer) CancelTopUp() func (ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req api.CancelTopUpRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		err = h.topUpService.CancelTopUp(ctx, req)
		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		h.RespondWithOK(ctx, nil)
	}
}
