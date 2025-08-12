package schema

type User struct {
	Name     string `json:"name" example:"Joe"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"joedoe@email.com"`
	Password string `json:"passwprd" example:"2r3c23rc3#@r32rs2"`
}
