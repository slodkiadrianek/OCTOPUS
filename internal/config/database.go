package config

import "database/sql"

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
