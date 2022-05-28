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
func DeleteTargetById(c *gin.Context) {
	var err error
	id := c.Param("id")

	target := domain.Target{}

	err = target.DeleteById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Target not deleted.")
		if err.Error() != "no rows in result set" {
			log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}

// Get all targets
func GetTargets(c *gin.Context) {

	var targets []domain.Target
	var err error
	var ilimit int64

	target := domain.Target{}

	name := c.DefaultQuery("name", "")
	limit := c.DefaultQuery("limit", "0")

	ilimit, err = strconv.ParseInt(limit, 10, 64)

	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid parameter limit.")

		log.WithFields(log.Fields{"limit": limit, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusBadRequest, serverError)
		return
	}

	targets, err = target.Read(util.Dbpool, name, ilimit)

	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Targets not found.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, targets)
}

// Get Target by Id
func GetTargetById(c *gin.Context) {
	id := c.Param("id")
	key := `"target: ` + id + `"`
	c.Header("ETag", key)

	ifnonematch := c.Request.Header.Get("If-None-Match")

	target := domain.Target{}

	if ifnonematch == key {
		c.IndentedJSON(http.StatusNotModified, nil)
		return
	}

	target, err := target.ReadById(util.Dbpool, id)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Target not found.")
		if err.Error() != "no rows in result set" { // Wrong id, does not exist
			log.WithFields(log.Fields{"id": id, "error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		}
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, target)
}

// Create new target
func PostTarget(c *gin.Context) {
	var newTarget domain.Target

	// Call BindJSON to bind the received JSON to newTarget.
	if err := c.BindJSON(&newTarget); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if newTarget.Id == 0 {
		log.WithFields(log.Fields{"newtarget": newTarget}).Debug("New target")
	}

	// Add the target to the database.
	_, err := newTarget.Write(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Newtarget not saved.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusInternalServerError, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newTarget)
}

// Update existing target
// See https://restfulapi.net/http-methods/
// Put only updates an existing target
//
func PutTargetById(c *gin.Context) {
	id := c.Param("id")
	var newTarget domain.Target

	// Call BindJSON to bind the received JSON to newTarget.
	if err := c.BindJSON(&newTarget); err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Error in json.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	if strconv.FormatInt(newTarget.Id, 10) != id {
		var serverError domain.ServerError = domain.GenerateServerError("Invalid identification of target, no modification.")
		log.WithFields(log.Fields{"clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusNotFound, serverError)
		return
	}

	// Update target in the database.
	_, err := newTarget.Update(util.Dbpool)
	if err != nil {
		var serverError domain.ServerError = domain.GenerateServerError("Target not updated.")

		log.WithFields(log.Fields{"error": err, "clientcode": serverError.Ticket}).Error(serverError.Message)
		c.IndentedJSON(http.StatusUnprocessableEntity, serverError)
		return
	}

	c.IndentedJSON(http.StatusOK, newTarget)
}
