package main

import (
	"context"
	"os"

	"github.com/bank/domain"
	"github.com/bank/test"

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

	test.AddAccounts(dbpool, maxnumber)
	accounts = test.ReadAccounts(dbpool, maxnumber)

	test.AddTargets(dbpool, maxnumber)
	test.AddTransactions(dbpool, maxnumber, accounts)
	log.Info("Closed  Bank")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
