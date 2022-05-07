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
	Read(dbpool *pgxpool.Pool, limit int) ([]Target, error)
	GetId() int64
	GetName() string
	GetDescription() string
	SetId(id int64)
	SetName(number string)
	SetDescription(description string)
}

func (target *Target) GetId() int64 {
	return target.id

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
	log.Debug("Write targets")

	log.WithFields(log.Fields{"id": target.id, "name": target.name, "description": target.description}).Debug("addTarget: Start addTarget")

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
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Debug("addTarget: insert target")
	}

	return lastInsertedId, nil
}

func (target *Target) Read(dbpool *pgxpool.Pool, limit int) ([]Target, error) {
	targets := []Target{}

	log.WithFields(log.Fields{"limit": limit}).Debug("Read target")

	rows, err := dbpool.Query(context.Background(), "SELECT * from target limit $1", limit)
	log.WithFields(log.Fields{"error": err}).Debug("Read account - after query")

	if err == nil {
		var index = 0

		for rows.Next() {
			log.WithFields(log.Fields{"index": index}).Debug("Read target - reading result")

			target := Target{}
			err := rows.Scan(&target.id, &target.name, &target.description)
			log.WithFields(log.Fields{"error": err}).Debug("Read account - reading result after scan error")
			if err == nil {
				targets = append(targets, target)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Debug("Read target - reading result error")
				return targets, err
			}
		}
		return targets, nil

	} else {
		log.WithFields(log.Fields{"error": err}).Debug("Read target - reading result error")
		return targets, err
	}
}
