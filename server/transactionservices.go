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
func DeleteTransactionById(c *gin.Context) {
	var err error
	id := c.Param("id")

	transaction := domain.Transaction{}

	err = transaction.DeleteById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Transaction not found, not deleted.")
		if err.Error() != "no rows in result set" {
			log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}

// Get all transactions
func GetTransactions(c *gin.Context) {

	var transactions []domain.Transaction
	var err error
	var ilimit int64

	from_account := c.DefaultQuery("from", "")
	to_account := c.DefaultQuery("to", "")
	limit := c.DefaultQuery("limit", "0")

	transaction := domain.Transaction{}

	ilimit, err = strconv.ParseInt(limit, 10, 64)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid parameter limit.")

		log.WithFields(log.Fields{"limit": limit, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusBadRequest, serverError)
		return
	}

	// retrieve known transactions
	transactions, err = transaction.Read(util.Dbpool, from_account, to_account, ilimit)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Transaction not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// convert accounts to json
	transactionstring, err := util.StrucToJsonString(transactions)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error converting transactions to json")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	// calculate hash and check if key already known in cache
	key := util.EtagHash(transactionstring)
	ifnonematch := c.Request.Header.Get("If-None-Match")
	log.WithFields(log.Fields{"If-None-Match": ifnonematch}).Trace("Before etag value")

	// return that value already present in client cache
	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	// return new value
	c.Header("Cache-Control", "max-age=30,  must-revalidate") // max-age in seconds
	c.Header("ETag", key)
	c.IndentedJSON(http.StatusOK, transactions)
}

// Get transaction by Id
func GetTransactionById(c *gin.Context) {
	id := c.Param("id")

	transaction := domain.Transaction{}

	// retrieve known account
	transaction, err := transaction.ReadById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Transaction not found.")
		if err.Error() != "no rows in result set" { // Wrong id, does not exist
			log.WithFields(log.Fields{"id": id, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// convert account to json
	transactionstring, err := util.StrucToJsonString(transaction)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error converting account to json")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	// calculate hash and check if key already known in cache
	key := util.EtagHash(transactionstring)
	ifnonematch := c.Request.Header.Get("If-None-Match")
	log.WithFields(log.Fields{"If-None-Match": ifnonematch}).Trace("Before etag value")

	// return that value already present in client cache
	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	// return new value
	c.Header("Cache-Control", "max-age=30,  must-revalidate") // max-age in seconds
	c.Header("ETag", key)
	c.IndentedJSON(http.StatusOK, transaction)
}

// Create new transaction
func PostTransaction(c *gin.Context) {
	var newTransaction domain.Transaction

	// Call BindJSON to bind the received JSON to newTransaction.
	if err := c.BindJSON(&newTransaction); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if newTransaction.Id == 0 {
		log.WithFields(log.Fields{"newtransaction": newTransaction}).Debug("New transaction")
	}

	// Add the transaction to the database.
	_, err := newTransaction.Write(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Newtransaction not saved.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newTransaction)
}

// Update existing account
// See https://restfulapi.net/http-methods/
// Put only updates an existing account
//
func PutTransactionById(c *gin.Context) {
	id := c.Param("id")
	var newTransaction domain.Transaction

	// Call BindJSON to bind the received JSON to newAccount.
	if err := c.BindJSON(&newTransaction); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if strconv.FormatInt(newTransaction.Id, 10) != id {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid identification of transaction, no modification.")
		log.WithFields(log.Fields{"clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// Update transaction in the database.
	_, err := newTransaction.Update(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Transaction not updated.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newTransaction)
}
