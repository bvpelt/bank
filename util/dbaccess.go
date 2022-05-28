package util

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	log "github.com/sirupsen/logrus"
)

var Dbpool *pgxpool.Pool

func Dbaccess() (*pgxpool.Pool, error) {
	var err error
	var msg string

	Dbpool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	msg = "Open database"
	log.WithFields(log.Fields{"pool": Dbpool}).Info(msg)
	CheckError(msg, err)

	if err != nil {
		return nil, err
	}

	err = Dbpool.Ping(context.Background())
	msg = "Check database connection"
	log.Info(msg)
	CheckError(msg, err)

	if err != nil {
		return nil, err
	}

	log.Info("Connected to database!")
	return Dbpool, err
}

func CheckError(msg string, err error) {
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error(msg)
		panic(err)
	}
}
