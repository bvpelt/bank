package server

import (
	"net/http"
	"time"

	"github.com/bank/domain"
	"github.com/bank/util"
	"github.com/gin-contrib/cors"
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
func jsonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

/*
func enableCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ETag")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		//c.Next()
	}
}
*/

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
	router.Use(jsonMiddleware())
	//router.Use(enableCors())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "ETag"},
		ExposeHeaders:    []string{"Content-Length", "ETag"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:4200"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(favicon.New("./resources/favicon.ico")) // set favicon middleware

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	return srv
}
