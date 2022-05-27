package server

import (
	"net/http"

	"github.com/bank/domain"
	"github.com/bank/util"
	"github.com/gin-gonic/gin"
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
	router.GET("/accounts", GetAccounts)
	router.GET("/accounts/:id", GetAccountById)
	router.PUT("/accounts/:id", PutAccountById)
	router.DELETE("/accounts/:id", DeleteAccountById)
	router.GET("/accounts/number/:number", GetAccountByNumber)
	router.POST("/accounts", PostAccount)
	router.GET("/pool", GetPool)
	router.Use(JSONMiddleware())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return srv
}
