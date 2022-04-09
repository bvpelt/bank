package main

import (
	"context"

	"fmt"
	"os"

	"github.com/jackc/pgx/v4"

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

var conn *pgx.Conn

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
	// urlExample := "postgres://testuser:12345@localhost:5432/bank"
	//psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	log.WithFields(log.Fields{"conn": conn}).
		Info("Open database")

		/*
			// open database
			db, err := sql.Open("postgres", psqlconn)
		*/
	CheckError(err)

	// close database
	defer func() {
		log.Info("Close database")
	}()

	/*
		// check db
		err = pgx.
		CheckError(err)
	*/
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
		id:          1,
		number:      "1234567890",
		description: "none",
	}

	id, err := addAccount(account)

	if err != nil {
		log.Error("addAccounts: Error during insert account")
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Added account")
	}
}

func addAccount(account Account) (int64, error) {

	log.WithFields(log.Fields{"id": account.id, "number": account.number, "description": account.description}).Info("addAccount: Start addAccount")

	pgCommandTag, err := conn.Exec(context.Background(), "INSERT INTO account (id, number, description) VALUES ($1, $2, $3)", account.id, account.number, account.description)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "account": account}).Error("addAccount: Error during insert account")
		return 0, fmt.Errorf("addAccount insert: %v", err)
	} else {
		log.WithFields(log.Fields{"pgCommandTag": pgCommandTag}).Info("addAccount: insert account")
	}

	/*
		id, err := result.LastInsertId()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("addAccount: Error during retrieve id from insert account")
			return 0, fmt.Errorf("addAccount get id: %v", err)
		}
	*/
	return 0, nil
}

func addTargets() {
	log.Info("Add targets")
}

func addTransactions() {
	log.Info("Add transactions")
}
