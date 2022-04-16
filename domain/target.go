package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Target struct {
	id          int64
	name        string
	description string
}

type ITarget interface {
	Write(dbpool *pgxpool.Pool) (int64, error)
}

func (target *Target) GetName() string {
	return target.name

}
func (target *Target) GetDescription() string {
	return target.description
}

func (target *Target) SetId(id int64) {
	target.id = id
}
func (target *Target) SetName(name string) {
	target.name = name
}

func (target *Target) SetDescription(description string) {
	target.description = description
}

func (target *Target) SetNameDescription(name string, description string) {
	target.name = name
	target.description = description
}

func (target *Target) Write(dbpool *pgxpool.Pool) (int64, error) {
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
