package DTO

type User struct {
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
}

func NewUser(email string, name string, surname string) *User {
	return &User{
		Email:   email,
		Name:    name,
		Surname: surname,
	}
}
