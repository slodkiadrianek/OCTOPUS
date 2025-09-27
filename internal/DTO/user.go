package DTO

type CreateUser struct {
	Name     string `json:"name" example:"Joe"`
	Surname  string `json:"surname" example:"Doe"`
	Email    string `json:"email" example:"joedoe@email.com"`
	Password string `json:"password" example:"2r3c23rc3#@r32rs2"`
}
type LoginUser struct {
	Email    string `json:"email" example:"adikurek@gmail.com"`
	Password string `json:"password" example:"zaqwekflas;h#&"`
}
type UpdateUser struct {
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
	Email   string `json:"email" example:"joedoe@email.com"`
}
type UserId struct {
	UserId string `json:"userId" example:"2"`
}
type ChangeUserPassword struct {
	CurrentPassword string `json:"currentPassword" example:"zaqw@Dekflas;h#&"`
	ConfirmPassword string `json:"confirmPassword" example:"zaqw@Dekflas;h#&"`
	NewPassword     string `json:"newPassword" example:"zaqw@Dekflas;h#&"`
}
type DeleteUser struct {
	Password string `json:"password" example:"zaqw@Dekflas;h#&"`
}
type UpdateUserNotifications struct {
	DiscordNotifications bool `json:"discordNotifications" example:"true"`
	SlackNotifications   bool `json:"slackNotifications" example:"true"`
	EmailNotifications   bool `json:"emailNotifications" example:"true"`
}
type LoggedUser struct {
	Id      int    `json:"id" example:"11"`
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
}

func NewLoggedUser(id int, email string, name string, surname string) *LoggedUser {
	return &LoggedUser{
		Id:      id,
		Email:   email,
		Name:    name,
		Surname: surname,
	}
}

func NewCreateUser(email string, name string, surname string) *CreateUser {
	return &CreateUser{
		Email:   email,
		Name:    name,
		Surname: surname,
	}
}
