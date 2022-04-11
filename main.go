package main

import (
	"context"

	"fmt"
	"os"

	//	"github.com/jackc/pgconn"
	//	"github.com/jackc/pgx"

	"github.com/jackc/pgx/v4/pgxpool"

	log "github.com/sirupsen/logrus"
)

type Account struct {
	id          int64
	number      string
	description string
}

type Target struct {
	id          int64
	name        string
	description string
}

type Transaction struct {
	id           int64
	from_account int64
	to_account   int64
	target       int64
	amount       int64
	description  string
}

//var conn *pgx.Conn
var dbpool *pgxpool.Pool

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)

}

func main() {
	log.Info("Started Bank")
	// export DATABASE_URL=postgres://testuser:12345@localhost:5432/bank

	var err error
	dbpool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	log.WithFields(log.Fields{"pool": dbpool}).
		Info("Open database")
	CheckError(err)

	err = dbpool.Ping(context.Background())
	CheckError(err)

	// close database
	defer func() {
		log.Info("Close database")
		dbpool.Close()
	}()

	log.Info("Connected to database!")

	addAccounts()
	addTargets()
	addTransactions()
	log.Info("Closed  Bank")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func addAccounts() {
	log.Info("Add accounts")
	account := Account{
		number:      "1234567890",
		description: "none",
	}

	id, err := addAccount(&account)
	fmt.Printf("%v %T\n", account, account)

	if err != nil {
		log.Error("addAccounts: Error during insert account")
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Added account")
	}
}

func addAccount(account *Account) (int64, error) {

	log.WithFields(log.Fields{"id": account.id, "number": account.number, "description": account.description}).Info("addAccount: Start addAccount")

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

func addTargets() {
	log.Info("Add targets")
	target := Target{
		name:        "Car",
		description: "none",
	}

	id, err := addTarget(&target)
	fmt.Printf("%v %T\n", target, target)

	if err != nil {
		log.Error("addTargets: Error during insert target")
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Added target")
	}

}

func addTarget(target *Target) (int64, error) {
	log.Info("Add targets")

	log.WithFields(log.Fields{"id": target.id, "name": target.name, "description": target.description}).Info("addTarget: Start addTarget")

	var err error
	var lastInsertedId int64 = 0

	if target.id != 0 {
		_, err = dbpool.Exec(context.Background(), "INSERT INTO target (id, name, description) VALUES ($1, $2, $3)", target.id, target.name, target.description)
		lastInsertedId = target.id
	} else {
		err = dbpool.QueryRow(context.Background(), "INSERT INTO target (name, description) VALUES ($1, $2) RETURNING id", target.name, target.description).Scan(&lastInsertedId)
		target.id = lastInsertedId
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "target": target}).Error("addTarget: Error during insert target")
		return 0, fmt.Errorf("addTarget insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Info("addTarget: insert target")
	}

	return lastInsertedId, nil
}

func addTransactions() {
	log.Info("Add transactions")
}
