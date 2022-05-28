package domain

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Target struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ITarget interface {
	DeleteById(dbpool *pgxpool.Pool, id string) error
	Read(dbpool *pgxpool.Pool, name string, limit int64) ([]Target, error)
	ReadById(dbpool *pgxpool.Pool, id string) (Target, error)
	Write(dbpool *pgxpool.Pool) (int64, error)

	GetDescription() string
	GetId() int64
	GetName() string

	SetDescription(description string)
	SetId(id int64)
	SetName(number string)
}

func (target *Target) DeleteById(dbpool *pgxpool.Pool, id string) error {
	var tar Target
	var err error

	// check if target exists
	rows := dbpool.QueryRow(context.Background(), "SELECT * from target where id = $1", id)

	err = rows.Scan(&tar.Id, &tar.Name, &tar.Description)

	if err == nil {
		_, err = dbpool.Exec(context.Background(), "DELETE from target where id = $1", id)
		log.WithFields(log.Fields{"error": err}).Debug("Delete target")
	}
	return err
}

func (target *Target) Read(dbpool *pgxpool.Pool, name string, limit int64) ([]Target, error) {
	var rows pgx.Rows
	var err error
	var query string = "SELECT * from target"
	var orderby string = " order by id desc"
	var where string = " where name='" + name + "'"

	targets := []Target{}

	if len(name) > 0 {
		query = query + where
	}

	query = query + orderby

	if limit > 0 {
		query = query + " limit $1"
		rows, err = dbpool.Query(context.Background(), query, limit)
	} else {
		rows, err = dbpool.Query(context.Background(), query)
	}

	if err == nil {
		var index = 0

		for rows.Next() {
			target := Target{}
			err := rows.Scan(&target.Id, &target.Name, &target.Description)

			if err == nil {
				targets = append(targets, target)
				index++
			} else {
				log.WithFields(log.Fields{"error": err}).Error("Read target - reading result error")
				return targets, err
			}
		}
		return targets, nil

	} else {
		if err.Error() != "no rows in result set" { // nothing found functional error
			log.WithFields(log.Fields{"error": err}).Error("Read target - reading result error")
		}
		return targets, err
	}
}

func (target *Target) ReadById(dbpool *pgxpool.Pool, id string) (Target, error) {
	var tar Target

	rows := dbpool.QueryRow(context.Background(), "SELECT * from target where id = $1", id)

	err := rows.Scan(&tar.Id, &tar.Name, &tar.Description)
	log.WithFields(log.Fields{"error": err, "target": target}).Trace("Read target - reading result after scan error")

	if err == nil {
		return tar, err
	} else {
		if err.Error() != "no rows in result set" { // wrong id, functional error
			log.WithFields(log.Fields{"id": id, "error": err}).Error("Read target - reading result error")
		}
		return tar, err
	}
}

func (target *Target) Update(dbpool *pgxpool.Pool) (int64, error) {

	var err error
	var lastInsertedId int64 = 0

	if target.Id != 0 {
		_, err = dbpool.Exec(context.Background(), "UPDATE target set name = $2, description = $3 where id = $1", target.Id, target.Name, target.Description)
		lastInsertedId = target.Id
	} else {
		return lastInsertedId, fmt.Errorf("identification for target is missing")
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err, "target": target}).Error("update target: Error during update target")
		return 0, fmt.Errorf("addTarget insert: %v", err)
	} else {
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Trace("update target: update target")
	}

	return lastInsertedId, nil
}

func (target *Target) Write(dbpool *pgxpool.Pool) (int64, error) {
	log.Debug("Write targets")

	log.WithFields(log.Fields{"id": target.Id, "name": target.Name, "description": target.Description}).Debug("addTarget: Start addTarget")

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
		log.WithFields(log.Fields{"lastInsertedId": lastInsertedId}).Debug("addTarget: insert target")
	}

	return lastInsertedId, nil
}

func (target *Target) GetId() int64 {
	return target.Id
}

func (target *Target) GetName() string {
	return target.Name
}

func (target *Target) GetDescription() string {
	return target.Description
}

func (target *Target) SetId(id int64) {
	target.Id = id
}

func (target *Target) SetName(name string) {
	target.Name = name
}

func (target *Target) SetDescription(description string) {
	target.Description = description
}

func (target *Target) SetNameDescription(name string, description string) {
	target.Name = name
	target.Description = description
}
