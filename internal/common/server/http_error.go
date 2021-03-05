package server

import (
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/common/errorz"
	"net/http"
)

func OK(ctx *gin.Context, resp interface{}) {
	ctx.JSON(http.StatusOK, resp)
}

func BadRequest(ctx *gin.Context, resp errorz.Response, err error) {
	httpRespondWithError(ctx, resp, err, "Bad Request", http.StatusBadRequest)
}

func NotFound(ctx *gin.Context, resp errorz.Response, err error) {
	httpRespondWithError(ctx, resp, err, "Not Found", http.StatusNotFound)
}

func InternalError(ctx *gin.Context, resp errorz.Response, err error) {
	httpRespondWithError(ctx, resp, err, "Internal Server Error", http.StatusInternalServerError)
}

func httpRespondWithError(ctx *gin.Context, resp errorz.Response, err error, logMSg string, status int) {
	// TODO: logging
	ctx.JSON(status, resp)
}
