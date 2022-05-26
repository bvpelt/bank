package server

import (
	"net/http"

	"github.com/bank/domain"
	"github.com/bank/util"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// get all accounts
func GetAccounts(c *gin.Context) {

	account := domain.Account{}
	var accounts []domain.Account

	accounts, err := account.Read(util.Dbpool, 0)
	if err == nil {
		c.IndentedJSON(http.StatusOK, accounts)
	} else {
		var serverError domain.ServerError = domain.GenerateServerError("Accounts not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
	}
}

// getAlbums responds with the list of all albums as JSON.
func GetAccountById(c *gin.Context) {
	id := c.Param("id")

	account := domain.Account{}

	account, err := account.ReadById(util.Dbpool, id)
	if err == nil {
		c.IndentedJSON(http.StatusOK, account)
	} else {
		var serverError domain.ServerError = domain.GenerateServerError("Account not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
	}
}

// getAlbums responds with the list of all albums as JSON.
func GetAccountByNumber(c *gin.Context) {
	number := c.Param("number")
	account := domain.Account{}
	var accounts []domain.Account

	accounts, err := account.ReadByNumber(util.Dbpool, number)
	if err == nil {
		if len(accounts) > 0 {
			c.IndentedJSON(http.StatusOK, accounts)
		} else {
			c.IndentedJSON(http.StatusNotFound, nil) // no technical error
		}
	} else {
		var serverError domain.ServerError = domain.GenerateServerError("Account not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
	}
}

// postAccount adds an account from JSON received in the request body.
func PostAccount(c *gin.Context) {
	var newAccount domain.Account

	// Call BindJSON to bind the received JSON to newAccount.
	if err := c.BindJSON(&newAccount); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if newAccount.Id == 0 {
		log.WithFields(log.Fields{"newacount": newAccount}).Debug("New accounts")
	}

	// Add the new account to the database.
	_, err := newAccount.Write(util.Dbpool)
	if err == nil {
		//	albums = append(albums, newAlbum)
		c.IndentedJSON(http.StatusCreated, newAccount)
	} else {
		var serverError domain.ServerError = domain.GenerateServerError("Newaccount not saved.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
	}
}

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
