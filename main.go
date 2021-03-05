package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opendigitalpay-io/open-balance/internal/common/server"
	"github.com/opendigitalpay-io/open-balance/internal/common/uid"
	"github.com/opendigitalpay-io/open-balance/internal/idem"
	"github.com/opendigitalpay-io/open-balance/internal/pay"
	"github.com/opendigitalpay-io/open-balance/internal/port"
	"github.com/opendigitalpay-io/open-balance/internal/storage"
	"github.com/opendigitalpay-io/open-balance/internal/topup"
	"github.com/opendigitalpay-io/open-balance/internal/transaction"
	"github.com/opendigitalpay-io/open-balance/internal/user"
	"net/http"
)

func main() {
	ctx := context.Background()
	repository, err := storage.NewRepository(ctx, &storage.Config{})
	if err != nil {
		panic(err)
	}

	uidGenerator, err := uid.NewGenerator(ctx)

	if err != nil {
		panic(err)
	}

	transactionService := transaction.NewService(repository, uidGenerator)
	idemService := idem.NewService(repository, uidGenerator)
	userService := user.NewService(repository, uidGenerator)
	topUpService := topup.NewService(repository, uidGenerator, transactionService, idemService)
	payService := pay.NewService(repository, uidGenerator, transactionService, idemService)

	server.RunHTTPServer(func(engine *gin.Engine) http.Handler {
		return port.HandlerFromMux(
			port.NewHTTPServer(userService, topUpService, payService, idemService, transactionService),
			engine,
		)
	})
}
