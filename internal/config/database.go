package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Db struct {
	DbConnection *sql.DB
}

func NewDb(databaseLink, driver string) (*Db, error) {
	dbConnection, err := sql.Open(driver, databaseLink)
	if err != nil {
		return &Db{}, err
	}

	err = dbConnection.Ping()
	if err != nil {
		return &Db{}, err
	}

	return &Db{
		DbConnection: dbConnection,
	}, nil
}
