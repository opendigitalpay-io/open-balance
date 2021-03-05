package port

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/port/api"
	"net/http"
)

func (h *HTTPServer) Start() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var idemReq api.IdemStartRequest
		err := ctx.ShouldBindJSON(&idemReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		println(idemReq.IdemKey)
		err = h.idemService.Start(ctx, idemReq.IdemKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, "Good")
	}
}

func (h *HTTPServer) End() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var idemReq api.IdemEndRequest
		err := ctx.ShouldBindJSON(&idemReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}
		println(idemReq.IdemKey)
		err = h.idemService.End(ctx, idemReq.IdemKey, idemReq.Response)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, "Good")
	}
}

func (h *HTTPServer) Test() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var idemReq api.IdemStartRequest
		err := ctx.ShouldBindJSON(&idemReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"err": err.Error(),
			})
			return
		}

		_, err = h.idemService.IdemExec(ctx, idemReq.IdemKey, func() (interface{}, error) {
			return "sample response", nil
		})

		if err != nil {
			h.RespondWithError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, "Good")
	}
}
