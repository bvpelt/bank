package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	Id           int64  `json:"id"`
	From_account int64  `json:"from"`
	To_account   int64  `json:"to"`
	Target       int64  `json:"target"`
	Amount       int64  `json:"amount"`
	Description  string `json:"description"`
}

type ITransaction interface {
	DeleteById(dbpool *pgxpool.Pool, id string) error
	Read(dbpool *pgxpool.Pool, from string, to string, limit int64) ([]Transaction, error)
	ReadById(dbpool *pgxpool.Pool) (Transaction, error)
	Update(dbpool *pgxpool.Pool) (int64, error)
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
	AddTransaction(dbpool *pgxpool.Pool) (int64, error)
}

func (transaction *Transaction) DeleteById(dbpool *pgxpool.Pool, id string) error {
	var trans Transaction
	var err error

	// check if transaction exists
	rows := dbpool.QueryRow(context.Background(), "SELECT * from transaction where id = $1", id)

	err = rows.Scan(&trans.Id, &trans.From_account, &trans.To_account, &trans.Target, &trans.Amount, &trans.Description)

	if err == nil {
		_, err = dbpool.Exec(context.Background(), "DELETE from transaction where id = $1", id)
		log.WithFields(log.Fields{"error": err}).Trace("Delete transaction")
	}
	return err
}

func (transaction *Transaction) Read(dbpool *pgxpool.Pool, from_account string, to_account string, limit int64) ([]Transaction, error) {
	var rows pgx.Rows
	var err error
	var query string = "SELECT * from transaction"
	var orderby string = " order by id desc"
	var where string = " where"
	var fromcrit string = " from_account='" + from_account + "'"
	var tocrit string = " to_account='" + to_account + "'"

	transactions := []Transaction{}

	if len(from_account) > 0 {
		query = query + where + fromcrit
		if len(to_account) > 0 {
			query = query + " and" + tocrit
		}

	} else {
		if len(to_account) > 0 {
			query = query + where + tocrit
		}
	}

	query = query + orderby

	if limit > 0 {
		query = query + " limit $1"
		rows, err = dbpool.Query(context.Background(), query, limit)
	} else {
		rows, err = dbpool.Query(context.Background(), query)
	}
	log.WithFields(log.Fields{"query": query}).Trace("Query to get all transactions")

	if err == nil {
		var index = 0

		for rows.Next() {
			transaction := Transaction{}
			err := rows.Scan(&transaction.Id, &transaction.From_account, &transaction.To_account, &transaction.Target, &transaction.Amount, &transaction.Description)

			if err == nil {
				transactions = append(transactions, transaction)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Error("Read transactions - reading result error")
				return transactions, err
			}
		}
		return transactions, nil
	} else {
		if err.Error() != "no rows in result set" { // nothing found functional error
			log.WithFields(log.Fields{"error": err}).Error("Read transaction - reading result error")
		}
		return transactions, err
	}
}

func (transaction *Transaction) ReadById(dbpool *pgxpool.Pool, id string) (Transaction, error) {
	var trans Transaction

	rows := dbpool.QueryRow(context.Background(), "SELECT * from transaction where id = $1", id)

	err := rows.Scan(&trans.Id, &trans.From_account, &trans.To_account, &trans.Target, &trans.Amount, &trans.Description)
	log.WithFields(log.Fields{"error": err, "transaction": trans}).Trace("Read transaction - reading result after scan error")

	if err == nil {
		return trans, err
	} else {
		if err.Error() != "no rows in result set" { // wrong id, functional error
			log.WithFields(log.Fields{"id": id, "error": err}).Error("Read transaction - reading result error")
		}
		return trans, err
	}
}

func (transaction *Transaction) Update(dbpool *pgxpool.Pool) (int64, error) {
	var err error
	var lastInsertedId int64 = 0

	if transaction.Id != 0 {
		_, err = dbpool.Exec(context.Background(), "UPDATE transaction set from_account = '$2', to_account = '$3', target = $4, amount = $5, description = '$3' where id = $1", transaction.Id, transaction.From_account, transaction.To_account, transaction.Target, transaction.Amount, transaction.Description)
		lastInsertedId = transaction.Id
	} else {
		return lastInsertedId, fmt.Errorf("identification for transaction is missing")
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err, "transaction": transaction}).Error("update transaction: Error during update transaction")
		return 0, fmt.Errorf("update Transaction insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Trace("update transaction: update transaction")
	}

	return lastInsertedId, nil
}

func (transaction *Transaction) Write(dbpool *pgxpool.Pool) (int64, error) {
	log.Debug("Write transaction")

	log.WithFields(log.Fields{"id": transaction.Id,
		"from_account": transaction.From_account,
		"to_account":   transaction.To_account,
		"target id":    transaction.Target,
		"amount":       transaction.Amount,
		"description":  transaction.Description,
	}).Debug("addTransaction: Start addTransaction")

	var err error
	var lastInsertedId int64 = 0

	if transaction.Id != 0 {
		_, err = dbpool.Exec(context.Background(),
			"INSERT INTO transaction (id, from_account, to_account, target, amount, description) VALUES ($1, $2, $3, $4, $5, $6)",
			transaction.Id, transaction.From_account, transaction.To_account, transaction.Target, transaction.Amount, transaction.Description)
		lastInsertedId = transaction.Id
	} else {
		err = dbpool.QueryRow(context.Background(),
			"INSERT INTO transaction (from_account, to_account, target, amount, description) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			transaction.From_account, transaction.To_account, transaction.Target, transaction.Amount, transaction.Description).Scan(&lastInsertedId)
		transaction.Id = lastInsertedId
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err, "transaction": transaction}).Error("addTransaction: Error during insert target")
		return 0, fmt.Errorf("addTaaddTransactionrget insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Debug("addTransaction: insert target")
	}

	return lastInsertedId, err
}

func (transaction *Transaction) AddTransaction(dbpool *pgxpool.Pool) (int64, error) {
	log.Debug("Add transaction")

	id, err := transaction.Write(dbpool)

	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (transaction *Transaction) GetId() int64 {
	return transaction.Id
}

func (transaction *Transaction) GetFromAccount() int64 {
	return transaction.From_account
}

func (transaction *Transaction) GetToAccount() int64 {
	return transaction.To_account
}

func (transaction *Transaction) GetTarget() int64 {
	return transaction.Target
}

func (transaction *Transaction) GetAmount() int64 {
	return transaction.Amount
}

func (transaction *Transaction) GetDescription() string {
	return transaction.Description
}

func (transaction *Transaction) SetId(id int64) {
	transaction.Id = id
}

func (transaction *Transaction) SetFromAccount(from_account int64) {
	transaction.From_account = from_account
}

func (transaction *Transaction) SetToAccount(to_account int64) {
	transaction.To_account = to_account
}

func (transaction *Transaction) SetTarget(target int64) {
	transaction.Target = target
}

func (transaction *Transaction) SetAmount(amount int64) {
	transaction.Amount = amount
}

func (transaction *Transaction) SetDescription(description string) {
	transaction.Description = description
}
