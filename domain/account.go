package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Account struct {
	Id          int64
	Number      string
	Description string
}

type IAccount interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
}

func (account *Account) Write(dbpool *pgxpool.Pool) (int64, error) {

	log.WithFields(log.Fields{"id": account.Id, "number": account.Number, "description": account.Description}).Info("addAccount: Start addAccount")

	var err error
	var lastInsertedId int64 = 0

	if account.Id != 0 {
		_, err = dbpool.Exec(context.Background(), "INSERT INTO account (id, number, description) VALUES ($1, $2, $3)", account.Id, account.Number, account.Description)
		lastInsertedId = account.Id
	} else {
		err = dbpool.QueryRow(context.Background(), "INSERT INTO account (number, description) VALUES ($1, $2) RETURNING id", account.Number, account.Description).Scan(&lastInsertedId)
		account.Id = lastInsertedId
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "account": account}).Error("addAccount: Error during insert account")
		return 0, fmt.Errorf("addAccount insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Info("addAccount: insert account")
	}

	return lastInsertedId, nil
}
