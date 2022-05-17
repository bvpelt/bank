package main

import (
	"context"
	"os"
	"time"

	"github.com/bank/test"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

var dbpool *pgxpool.Pool

func init() {
	//
	// Logging levels are: Trace, Debug, Info, Warning, Error, Fatal and Panic see https://github.com/Sirupsen/logrus
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		PadLevelText:    true,
		TimestampFormat: time.RFC3339,
	})
	/*
		customFormatter := new(logrus.TextFormatter)
		log.SetFormatter(customFormatter)
		customFormatter.TimestampFormat = "2016-03-28 15:04:05.000"
		customFormatter.FullTimestamp = true
	*/
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("Started Bank")
	// export DATABASE_URL=postgres://testuser:12345@localhost:5432/bank

	var err error

	dbpool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	log.WithFields(log.Fields{"pool": dbpool}).
		Debug("Open database")
	CheckError(err)

	err = dbpool.Ping(context.Background())
	CheckError(err)

	// always close database at program exit
	defer func() {
		log.Debug("Close database")
		dbpool.Close()
	}()

	log.Debug("Connected to database!")

	//test.DoTransactionTest(dbpool)
	test.Server()

	log.Debug("Closed  Bank")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
