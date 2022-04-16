package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	id           int64
	from_account int64
	to_account   int64
	target       int64
	amount       int64
	description  string
}

type ITransaction interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
	GetId() int64
	GetFromAccount() int64
	GetToAccount() int64
	GetTarget() int64
	GetAmount() int64
	SetId(id int64)
	SetFromAccount(from_account int64)
	SetToAccount(to_account int64)
	SetTarget(target int64)
	SetAmount(amount int64)
	SetDescription(description string)
}

func (transaction *Transaction) GetId() int64 {
	return transaction.id

}

func (transaction *Transaction) GetFromAccount() int64 {
	return transaction.from_account

}

func (transaction *Transaction) GetToAccount() int64 {
	return transaction.to_account

}

func (transaction *Transaction) GetTarget() int64 {
	return transaction.target

}

func (transaction *Transaction) GetAmount() int64 {
	return transaction.amount

}

func (transaction *Transaction) GetDescription() string {
	return transaction.description

}

func (transaction *Transaction) SetId(id int64) {
	transaction.id = id

}

func (transaction *Transaction) SetFromAccount(from_account int64) {
	transaction.from_account = from_account

}

func (transaction *Transaction) SetToAccount(to_account int64) {
	transaction.to_account = to_account

}
func (transaction *Transaction) SetTarget(target int64) {
	transaction.target = target

}
func (transaction *Transaction) SetAmount(amount int64) {
	transaction.amount = amount

}
func (transaction *Transaction) SetDescription(description string) {
	transaction.description = description

}

func (transaction *Transaction) Write(dbpool *pgxpool.Pool) (int64, error) {
	log.Info("Add targets")

	log.WithFields(log.Fields{"id": transaction.id,
		"from_account": transaction.from_account,
		"to_account":   transaction.to_account,
		"target id":    transaction.target,
		"amount":       transaction.amount,
		"description":  transaction.description,
	}).Info("addTarget: Start addTarget")

	var err error
	var lastInsertedId int64 = 0

	if transaction.id != 0 {
		_, err = dbpool.Exec(context.Background(),
			"INSERT INTO transaction (id, from_account, to_account, target, amount, description) VALUES ($1, $2, $3, $4, $5, $6)",
			transaction.id, transaction.from_account, transaction.to_account, transaction.target, transaction.amount, transaction.description)
		lastInsertedId = transaction.id
	} else {
		err = dbpool.QueryRow(context.Background(),
			"INSERT INTO transaction (from_account, to_account, target, amount, description) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			transaction.from_account, transaction.to_account, transaction.target, transaction.amount, transaction.description).Scan(&lastInsertedId)
		transaction.id = lastInsertedId
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "transaction": transaction}).Error("addTransaction: Error during insert target")
		return 0, fmt.Errorf("addTarget insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Info("addTransaction: insert target")
	}

	return lastInsertedId, nil
}
