package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bank/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Transaction struct {
	Id           int64
	From_account int64
	To_account   int64
	Target       int64
	Amount       int64
	Description  string
}

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
	account := domain.Account{
		Number:      "1234567890",
		Description: "none",
	}

	//id, err := domain.AddAccount(dbpool, &account)
	id, err := account.Write(dbpool)
	fmt.Printf("%v %T\n", account, account)

	if err != nil {
		log.Error("addAccounts: Error during insert account")
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Added account")
	}
}

func addTargets() {
	log.Info("Add targets")
	target := domain.Target{
		Name:        "Car",
		Description: "none",
	}

	id, err := target.Write(dbpool)
	fmt.Printf("%v %T\n", target, target)

	if err != nil {
		log.Error("addTargets: Error during insert target")
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Added target")
	}

}

func addTransactions() {
	log.Info("Add transactions")
}
