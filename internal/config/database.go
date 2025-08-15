package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Db struct {
	DbConnection *sql.DB
}

func NewDb(databaseLink string) *Db {
	dbConnection, err := sql.Open("postgres", databaseLink)
	if err != nil {
		panic(err)
	}
	return &Db{
		DbConnection: dbConnection,
	}
}
