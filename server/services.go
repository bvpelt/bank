package server

import (
	"net/http"

	"github.com/bank/domain"
	"github.com/bank/util"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// get pool information
func GetPool(c *gin.Context) {
	stat := util.Dbpool.Stat()

	var status domain.DbpoolStat
	status.AcquireConns = stat.AcquiredConns()
	status.AcquireCount = stat.AcquireCount()
	status.AcquireDuration = stat.AcquireDuration()
	status.ConstructingConns = stat.ConstructingConns()
	status.EmptyAcquireCount = stat.EmptyAcquireCount()
	status.IdleConns = stat.IdleConns()
	status.MaxConns = stat.MaxConns()
	status.TotalConns = stat.TotalConns()

	c.IndentedJSON(http.StatusOK, status)
}

// use contenttype application/json for all services
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func StartServer() *http.Server {

	router := gin.Default()
	router.DELETE("/accounts/:id", DeleteAccountById)
	router.GET("/accounts", GetAccounts)
	router.GET("/accounts/:id", GetAccountById)
	router.POST("/accounts", PostAccount)
	router.PUT("/accounts/:id", PutAccountById)
	router.GET("/accounts/search/:term", SearchAccounts)

	router.DELETE("/targets/:id", DeleteTargetById)
	router.GET("/targets", GetTargets)
	router.GET("/targets/:id", GetTargetById)
	router.POST("/targets", PostTarget)
	router.PUT("/targets/:id", PutTargetById)

	router.DELETE("/transactions/:id", DeleteTargetById)
	router.GET("/transactions", GetTransactions)
	router.GET("/transactions/:id", GetTransactionById)
	router.POST("/transactions", PostTransaction)
	router.PUT("/transactions/:id", PutTransactionById)

	router.GET("/pool", GetPool)
	router.Use(JSONMiddleware())
	router.Use(favicon.New("./resources/favicon.ico")) // set favicon middleware

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return srv
}
