package DTO

type CreateUser struct {
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
}

type LoggedUser struct {
	Id      int    `json:"id" example:"11"`
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
}

func NewUser(email string, name string, surname string) *CreateUser {
	return &CreateUser{
		Email:   email,
		Name:    name,
		Surname: surname,
	}
}
