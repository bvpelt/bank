package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/bank/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

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
	var accounts []domain.Account

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

	const maxnumber = 20

	addAccounts(maxnumber)
	accounts = readAccounts(maxnumber)

	addTargets(maxnumber)
	addTransactions(maxnumber, accounts)
	log.Info("Closed  Bank")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func addAccounts(maxnumber int) {
	log.Info("Add accounts")
	account := domain.Account{}

	account.SetNumber("1234567890")
	account.SetDescription("none")
	for i := 0; i < maxnumber; i++ {
		str := strconv.Itoa(i)
		account.SetId(0)
		account.SetNumber(str)
		id, err := account.Write(dbpool)
		fmt.Printf("%v %T\n", account, account)

		if err != nil {
			log.Error("addAccounts: Error during insert account")
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Added account")
		}
	}
}

func readAccounts(maxnumber int) []domain.Account {
	var accounts []domain.Account
	var account domain.Account

	accounts, err := account.Read(dbpool, maxnumber)

	if err == nil {
		for _, account = range accounts {
			fmt.Println("Found account: %v\n", account)
		}
	}

	return accounts
}

func addTargets(maxnumber int) {
	log.Info("Add targets")
	target := domain.Target{}
	target.SetName("Car")
	target.SetDescription("none")

	for i := 0; i < maxnumber; i++ {
		str := "name: " + strconv.Itoa(i)
		target.SetId(0)
		target.SetName(str)

		id, err := target.Write(dbpool)
		fmt.Printf("%v %T\n", target, target)

		if err != nil {
			log.Error("addTargets: Error during insert target")
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Added target")
		}
	}
}

func addTransactions(maxnumber int, accounts []domain.Account) {
	log.Info("Add transactions")

	transaction := domain.Transaction{}
	transaction.SetTarget(1)

	var i int
	for i = 0; i < maxnumber; i++ {
		str := "transaction: " + strconv.Itoa(i)
		transaction.SetId(0)
		transaction.SetDescription(str)
		transaction.SetFromAccount(accounts[i].GetId())
		transaction.SetToAccount(accounts[maxnumber-(i+1)].GetId())
		transaction.SetAmount(int64(rand.Intn(25000)))

		id, err := transaction.Write(dbpool)
		fmt.Printf("%v %T\n", transaction, transaction)

		if err != nil {
			log.Error("addTransactions: Error during insert transaction")
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Added transaction")
		}
	}
}
