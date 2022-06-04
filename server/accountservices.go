package server

import (
	"net/http"
	"strconv"

	"github.com/bank/domain"
	"github.com/bank/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Delete Accounts by Id
func DeleteAccountById(c *gin.Context) {
	var err error
	id := c.Param("id")

	account := domain.Account{}

	err = account.DeleteById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Account not found, not deleted.")
		if err.Error() != "no rows in result set" {
			log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}

// Get all accounts
func GetAccounts(c *gin.Context) {

	var accounts []domain.Account
	var err error
	var ilimit int64

	number := c.DefaultQuery("number", "")
	limit := c.DefaultQuery("limit", "0")

	account := domain.Account{}

	ilimit, err = strconv.ParseInt(limit, 10, 64)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid parameter limit.")

		log.WithFields(log.Fields{"limit": limit, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusBadRequest, serverError)
		return
	}

	// retrieve known accounts
	accounts, err = account.Read(util.Dbpool, number, ilimit)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Accounts not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// convert accounts to json
	accountstring, err := util.StrucToJsonString(accounts)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error converting accounts to json")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	// calculate hash and check if key already known in cache
	key := util.EtagHash(accountstring)
	ifnonematch := c.Request.Header.Get("If-None-Match")
	log.WithFields(log.Fields{"If-None-Match": ifnonematch}).Trace("Before etag value")

	// return that value already present in client cache
	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	// return new value
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "max-age=30,  must-revalidate") // max-age in seconds
	c.Header("ETag", key)
	c.IndentedJSON(http.StatusOK, accounts)
}

// Get Account by Id
func GetAccountById(c *gin.Context) {
	id := c.Param("id")

	account := domain.Account{}

	// retrieve known account
	account, err := account.ReadById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Account not found.")
		if err.Error() != "no rows in result set" { // Wrong id, does not exist
			log.WithFields(log.Fields{"id": id, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// convert account to json
	accountstring, err := util.StrucToJsonString(account)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error converting account to json")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	// calculate hash and check if key already known in cache
	key := util.EtagHash(accountstring)
	ifnonematch := c.Request.Header.Get("If-None-Match")
	log.WithFields(log.Fields{"If-None-Match": ifnonematch}).Trace("Before etag value")

	// return that value already present in client cache
	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	// return new value
	c.Header("Cache-Control", "max-age=30,  must-revalidate") // max-age in seconds
	c.Header("ETag", key)
	c.IndentedJSON(http.StatusOK, account)
}

// Create new account
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
		log.WithFields(log.Fields{"newacount": newAccount}).Debug("New account")
	}

	// Add the account to the database.
	_, err := newAccount.Write(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Newaccount not saved.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newAccount)
}

// Update existing account
// See https://restfulapi.net/http-methods/
// Put only updates an existing account
//
func PutAccountById(c *gin.Context) {
	id := c.Param("id")
	var newAccount domain.Account

	// Call BindJSON to bind the received JSON to newAccount.
	if err := c.BindJSON(&newAccount); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if strconv.FormatInt(newAccount.Id, 10) != id {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid identification of account, no modification.")
		log.WithFields(log.Fields{"clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// Update account in the database.
	_, err := newAccount.Update(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Account not updated.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newAccount)
}
