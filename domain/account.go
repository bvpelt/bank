package domain

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Account struct {
	Id          int64  `json:"id"`
	Number      string `json:"number"`
	Description string `json:"description"`
}

type IAccount interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
	Read(dbpool *pgxpool.Pool, limit int) ([]Account, error)
	ReadById(dbpool *pgxpool.Pool) (Account, error)
	GetId() int64
	GetNumber() string
	GetDescription() string
	SetId(id int64)
	SetNumber(number string)
	SetDescription(description string)
	GetAccounts(c *gin.Context, dbpool *pgxpool.Pool)
}

func (account *Account) GetId() int64 {
	return account.Id
}

func (account *Account) GetNumber() string {
	return account.Number
}

func (account *Account) GetDescription() string {
	return account.Description
}

func (account *Account) SetId(id int64) {
	account.Id = id
}

func (account *Account) SetNumber(number string) {
	account.Number = number
}

func (account *Account) SetDescription(description string) {
	account.Description = description
}

func (account *Account) SetNumberDescription(number string, description string) {
	account.Number = number
	account.Description = description
}

func (account *Account) Write(dbpool *pgxpool.Pool) (int64, error) {

	log.WithFields(log.Fields{"id": account.Id, "number": account.Number, "description": account.Description}).Debug("Write account")

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
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Debug("addAccount: insert account")
	}

	return lastInsertedId, nil
}

func (account *Account) Read(dbpool *pgxpool.Pool, limit int) ([]Account, error) {
	accounts := []Account{}

	//log.WithFields(log.Fields{"limit": limit}).Debug("Read account")

	var query string = "SELECT * from account order by id desc"

	var rows pgx.Rows
	var err error

	if limit > 0 {
		query = query + " limit $1"
		rows, err = dbpool.Query(context.Background(), query, limit)
	} else {
		rows, err = dbpool.Query(context.Background(), query)

	}

	//rows, err = dbpool.Query(context.Background(), "SELECT * from account order by id limit $1", limit)
	//log.WithFields(log.Fields{"error": err}).Debug("Read account - after query")

	if err == nil {
		var index = 0

		for rows.Next() {
			//log.WithFields(log.Fields{"index": index}).Debug("Read account - reading result")

			account := Account{}
			err := rows.Scan(&account.Id, &account.Number, &account.Description)
			//log.WithFields(log.Fields{"error": err, "account": account}).Debug("Read account - reading result after scan error")
			if err == nil {
				accounts = append(accounts, account)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result error")
				return accounts, err
			}
		}
		return accounts, nil
	} else {
		log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result error")
		return accounts, err
	}
}

func (account *Account) ReadById(dbpool *pgxpool.Pool, id string) (Account, error) {
	var acc Account
	log.Debug("Read account by id")

	rows := dbpool.QueryRow(context.Background(), "SELECT * from account where id = $1", id)
	log.Debug("Read account by id - after query")

	err := rows.Scan(&acc.Id, &acc.Number, &acc.Description)
	log.WithFields(log.Fields{"error": err, "account": account}).Debug("Read account - reading result after scan error")
	if err == nil {
		return acc, err
	} else {
		log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result error")
		return acc, err
	}
}

func (account *Account) ReadByNumber(dbpool *pgxpool.Pool, number string) ([]Account, error) {
	accounts := []Account{}
	log.Debug("Read account by number")

	rows, err := dbpool.Query(context.Background(), "SELECT * from account where number = $1 order by number", number)
	log.Debug("Read account by id - after query")

	if err == nil {
		var index = 0

		for rows.Next() {
			account := Account{}
			err := rows.Scan(&account.Id, &account.Number, &account.Description)

			if err == nil {
				accounts = append(accounts, account)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result error")
				return accounts, err
			}
		}
		return accounts, nil
	} else {
		log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result error")
		return accounts, err
	}

}
