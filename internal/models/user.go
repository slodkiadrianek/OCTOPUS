package models

type User struct {
	Id       int    `json:"id" sql:"id" example:"1"`
	Email    string `json:"email" sql:"email" example:"joedoe@email.com"`
	Name     string `json:"name" sql:"name" example:"Joe"`
	Surname  string `json:"surname" sql:"surname" example:"Doe"`
	Password string `json:"password" example:"fsdf2332@!32"`
}

// var query string = "CREATE TABLE IF NOT EXISTS users (" +
//	"id INT PRIMARY KEY AUTOINCREMENT," +
//	"email VARCHAR(128) UNIQUE," +
//	"name VARCHAR(64)," +
//	"surname VARCHAR(64)," +
//	"role VARCHAR(64)," +
//	")"
