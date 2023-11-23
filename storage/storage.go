package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DidierWatson/go-db/pkg/product"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

type Driver string

const (
	MySQL    Driver = "MYSQL"
	Postgres Driver = "POSTGRES"
)

func New(d Driver) {
	switch d {
	case MySQL:
		newMySQLDB()
	case Postgres:
		newPostgresDB()
	}
}

func newPostgresDB() {
	once.Do(func() {
		var err error
		connStr := "postgres://postgres:@localhost:5432/test?sslmode=disable"
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("can't ping: %v", err)
		}
		fmt.Println("Conectado a postgre")
	})

}

func newMySQLDB() {
	once.Do(func() {
		var err error
		connStr := "root:Mulsanne.3@tcp(localhost:3306)/goMysql?parseTime=true"
		db, err = sql.Open("mysql", connStr)
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("can't ping: %v", err)
		}
		fmt.Println("Conectado a Mysql")
	})

}

// Pool returns unique instance
func Pool() *sql.DB {
	return db
}

func stringToNull(s string) sql.NullString {
	null := sql.NullString{String: s}
	if null.String != "" {
		null.Valid = true
	}
	return null
}

func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRowProduct(s scanner) (*product.Model, error) {
	m := &product.Model{}
	observationNull := sql.NullString{}
	updatedAtNull := sql.NullTime{}
	err := s.Scan(
		&m.ID,
		&m.Name,
		&observationNull,
		&m.Price,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &product.Model{}, err
	}
	m.Observations = observationNull.String
	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}

func DAOProduct(driver Driver) (product.Storage, error) {
	switch driver {
	case Postgres:
		return newpsqlProduct(db), nil
	case MySQL:
		return newMySQLProduct(db), nil
	default:
		return nil, fmt.Errorf("Driver not implemented")
	}
}
