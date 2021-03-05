package port

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opendigitalpay-io/open-balance/internal/common/errorz"
	"github.com/opendigitalpay-io/open-balance/internal/common/server"
	"github.com/opendigitalpay-io/open-balance/internal/domain"
	"github.com/opendigitalpay-io/open-balance/internal/idem"
	"github.com/opendigitalpay-io/open-balance/internal/pay"
	"github.com/opendigitalpay-io/open-balance/internal/storage"
	"github.com/opendigitalpay-io/open-balance/internal/topup"
	"github.com/opendigitalpay-io/open-balance/internal/transaction"
	"github.com/opendigitalpay-io/open-balance/internal/user"
)

type HTTPServer struct {
	userService        user.Service // *service
	topUpService       topup.Service
	payService         pay.Service
	idemService        idem.Service
	transactionService transaction.Service // FixMe: Temporary for demo, delete later
}

func NewHTTPServer(userService user.Service, topUpService topup.Service, payService pay.Service, idemService idem.Service, transactionService transaction.Service) *HTTPServer {
	return &HTTPServer{
		userService:        userService,
		topUpService:       topUpService,
		payService:         payService,
		idemService:        idemService,
		transactionService: transactionService, // FixMe: Temporary for demo, delete later
	}
}

func (*HTTPServer) RespondWithOK(ctx *gin.Context, resp interface{}) {
	server.OK(ctx, resp)
}

func (*HTTPServer) RespondWithError(ctx *gin.Context, err error) {
	var ves validator.ValidationErrors
	if errors.As(err, &ves) {
		server.BadRequest(ctx, errorz.NewValidationError(ves), err)
		return
	}

	var syne *json.SyntaxError
	if errors.As(err, &syne) {
		server.BadRequest(ctx, errorz.NewInvalidJSONError(syne), err)
		return
	}

	var nfe storage.NotFoundError
	if errors.As(err, &nfe) {
		server.NotFound(ctx, errorz.NewNotFoundError(nfe), err)
		return
	}

	var dee storage.DuplicatedEntryError
	if errors.As(err, &dee) {
		server.BadRequest(ctx, errorz.NewInvalidValueError(dee), err)
		return
	}

	var ideme domain.IdemError
	if errors.As(err, &ideme) {
		server.BadRequest(ctx, errorz.NewIdemKeyError(err), err)
		return
	}

	var txne domain.TransactionError
	if errors.As(err, &txne) {
		server.BadRequest(ctx, errorz.NewTransactionError(err), err)
		return
	}

	// fallback error handling
	server.InternalError(ctx, errorz.NewInternalError(err), err)
}
