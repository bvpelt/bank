package domain

import (
	"context"
	"fmt"

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
	DeleteById(dbpool *pgxpool.Pool, id string)
	Read(dbpool *pgxpool.Pool, number string, limit int64) ([]Account, error)
	ReadById(dbpool *pgxpool.Pool) (Account, error)
	ReadByNumber(dbpool *pgxpool.Pool, number string) ([]Account, error)
	Update(dbpool *pgxpool.Pool) (int64, error)
	Write(dbpool *pgxpool.Pool) (int64, error)
	GetDescription() string
	GetId() int64
	GetNumber() string
	SetDescription(description string)
	SetId(id int64)
	SetNumber(number string)
	SetNumberDescription(number string, description string)
}

func (account *Account) DeleteById(dbpool *pgxpool.Pool, id string) error {

	log.Debug("Delete account by id")

	_, err := dbpool.Exec(context.Background(), "DELETE from account where id = $1", id)
	log.WithFields(log.Fields{"error": err}).Trace("Delete account")
	return err
}

func (account *Account) Read(dbpool *pgxpool.Pool, number string, limit int64) ([]Account, error) {
	var rows pgx.Rows
	var err error
	var query string = "SELECT * from account"
	var orderby string = " order by id desc"
	var where string = " where number='" + number + "'"

	accounts := []Account{}

	if len(number) > 0 {
		query = query + where
	}

	query = query + orderby

	if limit > 0 {
		query = query + " limit $1"
		rows, err = dbpool.Query(context.Background(), query, limit)
	} else {
		rows, err = dbpool.Query(context.Background(), query)

	}

	if err == nil {
		var index = 0

		for rows.Next() {
			account := Account{}
			err := rows.Scan(&account.Id, &account.Number, &account.Description)

			if err == nil {
				accounts = append(accounts, account)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Error("Read account - reading result error")
				return accounts, err
			}
		}
		return accounts, nil
	} else {
		log.WithFields(log.Fields{"error": err}).Error("Read account - reading result error")
		return accounts, err
	}
}

func (account *Account) ReadById(dbpool *pgxpool.Pool, id string) (Account, error) {
	var acc Account

	rows := dbpool.QueryRow(context.Background(), "SELECT * from account where id = $1", id)

	err := rows.Scan(&acc.Id, &acc.Number, &acc.Description)
	log.WithFields(log.Fields{"error": err, "account": account}).Trace("Read account - reading result after scan error")

	if err == nil {
		return acc, err
	} else {
		if err.Error() != "no rows in result set" { // wrong id, functional error
			log.WithFields(log.Fields{"id": id, "error": err}).Error("Read account - reading result error")
		}
		return acc, err
	}
}

func (account *Account) ReadByNumber(dbpool *pgxpool.Pool, number string) ([]Account, error) {
	accounts := []Account{}

	rows, err := dbpool.Query(context.Background(), "SELECT * from account where number = $1 order by number", number)

	if err == nil {
		var index = 0

		for rows.Next() {
			account := Account{}
			err := rows.Scan(&account.Id, &account.Number, &account.Description)

			if err == nil {
				accounts = append(accounts, account)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Error("Read account - reading result error")
				return accounts, err
			}
		}
		return accounts, nil
	} else {
		if err.Error() != "no rows in result set" { // wrong number, functional error
			log.WithFields(log.Fields{"number": number, "error": err}).Error("Read account - reading result error")
		}
		return accounts, err
	}
}

func (account *Account) Update(dbpool *pgxpool.Pool) (int64, error) {

	var err error
	var lastInsertedId int64 = 0

	if account.Id != 0 {
		_, err = dbpool.Exec(context.Background(), "UPDATE account set number = $2, description = $3 where id = $1", account.Id, account.Number, account.Description)
		lastInsertedId = account.Id
	} else {
		return lastInsertedId, fmt.Errorf("Identification for account is missing")
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err, "account": account}).Error("addAccount: Error during insert account")
		return 0, fmt.Errorf("addAccount insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Trace("addAccount: insert account")
	}

	return lastInsertedId, nil
}

func (account *Account) Write(dbpool *pgxpool.Pool) (int64, error) {

	log.WithFields(log.Fields{"id": account.Id, "number": account.Number, "description": account.Description}).Trace("Write account")

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
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Trace("addAccount: insert account")
	}

	return lastInsertedId, nil
}

func (account *Account) GetDescription() string {
	return account.Description
}

func (account *Account) GetId() int64 {
	return account.Id
}

func (account *Account) GetNumber() string {
	return account.Number
}

func (account *Account) SetDescription(description string) {
	account.Description = description
}

func (account *Account) SetId(id int64) {
	account.Id = id
}

func (account *Account) SetNumber(number string) {
	account.Number = number
}

func (account *Account) SetNumberDescription(number string, description string) {
	account.Number = number
	account.Description = description
}
