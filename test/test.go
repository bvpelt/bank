package test

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/bank/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func AddAccounts(dbpool *pgxpool.Pool, maxnumber int, wg *sync.WaitGroup) {
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
	wg.Done()
}

func ReadAccounts(dbpool *pgxpool.Pool, maxnumber int) []domain.Account {
	var accounts []domain.Account
	var account domain.Account

	accounts, err := account.Read(dbpool, maxnumber)

	if err == nil {
		for _, account = range accounts {
			log.WithFields(log.Fields{"id": account.GetId(), "number": account.GetNumber(), "description": account.GetDescription()}).Info("found account")
		}
	}

	return accounts
}

func AddTargets(dbpool *pgxpool.Pool, maxnumber int, wg *sync.WaitGroup) {
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
			log.WithFields(log.Fields{"error": err}).Error("addTargets: Error during insert target")
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Added target")
		}
	}
	wg.Done()
}

func AddTransactions(dbpool *pgxpool.Pool, maxnumber int, accounts []domain.Account) {
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
