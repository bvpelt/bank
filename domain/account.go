package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Account struct {
	id          int64
	number      string
	description string
}

type IAccount interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
	Read(dbpool *pgxpool.Pool, limit int) ([]Account, error)
	GetId() int64
	GetNumber() string
	GetDescription() string
	SetId(id int64)
	SetNumber(number string)
	SetDescription(description string)
}

func (account *Account) GetId() int64 {
	return account.id

}

func (account *Account) GetNumber() string {
	return account.number

}

func (account *Account) GetDescription() string {
	return account.description
}

func (account *Account) SetId(id int64) {
	account.id = id
}
func (account *Account) SetNumber(number string) {
	account.number = number
}

func (account *Account) SetDescription(description string) {
	account.description = description
}

func (account *Account) SetNumberDescription(number string, description string) {
	account.number = number
	account.description = description
}

func (account *Account) Write(dbpool *pgxpool.Pool) (int64, error) {

	log.WithFields(log.Fields{"id": account.id, "number": account.number, "description": account.description}).Info("Write account")

	var err error
	var lastInsertedId int64 = 0

	if account.id != 0 {
		_, err = dbpool.Exec(context.Background(), "INSERT INTO account (id, number, description) VALUES ($1, $2, $3)", account.id, account.number, account.description)
		lastInsertedId = account.id
	} else {
		err = dbpool.QueryRow(context.Background(), "INSERT INTO account (number, description) VALUES ($1, $2) RETURNING id", account.number, account.description).Scan(&lastInsertedId)
		account.id = lastInsertedId
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "account": account}).Error("addAccount: Error during insert account")
		return 0, fmt.Errorf("addAccount insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Info("addAccount: insert account")
	}

	return lastInsertedId, nil
}

func (account *Account) Read(dbpool *pgxpool.Pool, limit int) ([]Account, error) {
	accounts := []Account{}

	log.WithFields(log.Fields{"limit": limit}).Info("Read account")

	rows, err := dbpool.Query(context.Background(), "SELECT * from account limit $1", limit)
	log.WithFields(log.Fields{"error": err}).Info("Read account - after query")

	if err == nil {
		var index = 0

		for rows.Next() {
			log.WithFields(log.Fields{"index": index}).Info("Read account - reading result")

			account := Account{}
			err := rows.Scan(&account.id, &account.number, &account.description)
			log.WithFields(log.Fields{"error": err}).Info("Read account - reading result after scan error")
			if err == nil {

				accounts = append(accounts, account)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Info("Read account - reading result error")
				return accounts, err
			}
		}
		return accounts, nil

	} else {
		log.WithFields(log.Fields{"error": err}).Info("Read account - reading result error")
		return accounts, err
	}

}
