package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Target struct {
	Id          int64
	Name        string
	Description string
}

type ITarget interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
}

func (target *Target) Write(dbpool *pgxpool.Pool) (int64, error) {
	log.Info("Add targets")

	log.WithFields(log.Fields{"id": target.Id, "name": target.Name, "description": target.Description}).Info("addTarget: Start addTarget")

	var err error
	var lastInsertedId int64 = 0

	if target.Id != 0 {
		_, err = dbpool.Exec(context.Background(), "INSERT INTO target (id, name, description) VALUES ($1, $2, $3)", target.Id, target.Name, target.Description)
		lastInsertedId = target.Id
	} else {
		err = dbpool.QueryRow(context.Background(), "INSERT INTO target (name, description) VALUES ($1, $2) RETURNING id", target.Name, target.Description).Scan(&lastInsertedId)
		target.Id = lastInsertedId
	}
	if err != nil {
		log.WithFields(log.Fields{"error": err, "target": target}).Error("addTarget: Error during insert target")
		return 0, fmt.Errorf("addTarget insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Info("addTarget: insert target")
	}

	return lastInsertedId, nil
}
