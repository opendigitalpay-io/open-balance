package port

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerInterface interface {
	// GET /health
	GetHealthStatus() func(*gin.Context)

	//  GET /v1/user/:id
	GetUserByID() func(*gin.Context)
	//  POST /v1/user
	AddUser() func(*gin.Context)
	//  POST /v1/user/:id
	UpdateUser() func(*gin.Context)

	// POST /v1/topup/try
	TryTopUp() func(*gin.Context)
	// POST /v1/topup/commit
	CommitTopUp() func(*gin.Context)
	// POST /v1/topup/cancel
	CancelTopUp() func(*gin.Context)

	// POST /v1/pay/try
	TryPay() func(*gin.Context)
	// POST /v1/pay/commit
	CommitPay() func(*gin.Context)
	// POST /v1/pay/cancel
	CancelPay() func(*gin.Context)

	// POST /v1/idem/start
	Start() func(*gin.Context)
	// POST /v1/idem/end
	End() func(*gin.Context)
	// POST /v1/idem/test
	Test() func(*gin.Context)

	// FixMe: Temporary for demo, delete later
	// GET /v1/transaction/:userId
	GetTransactionByUserID() func(*gin.Context)
}

func HandlerFromMux(si ServerInterface, e *gin.Engine) http.Handler {
	// Healthz
	e.GET("/health", si.GetHealthStatus())

	// v1
	v1 := e.Group("/v1")
	{
		// user
		u := v1.Group("/user")
		{
			u.GET("/:id", si.GetUserByID())
			u.POST("", si.AddUser())
			u.POST("/:id", si.UpdateUser())
		}

		// top-up
		t := v1.Group("/topup")
		{
			t.POST("/try", si.TryTopUp())
			t.POST("/commit", si.CommitTopUp())
			t.POST("/cancel", si.CancelTopUp())
		}

		// pay
		p := v1.Group("/pay")
		{
			p.POST("/try", si.TryPay())
			p.POST("/commit", si.CommitPay())
			p.POST("/cancel", si.CancelPay())
		}

		// idem
		idem := v1.Group("/idem")
		{
			idem.POST("/start", si.Start())
			idem.POST("/end", si.End())
			idem.POST("/test", si.Test())
		}

		// FixMe: Temporary for demo, delete later
		// transaction
		txn := v1.Group("/transaction")
		{
			txn.GET("/:userId", si.GetTransactionByUserID())
		}
	}

	return e
}
