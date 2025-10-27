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
type UserID struct {
	UserID string `json:"userID" example:"2"`
}
type ChangeUserPassword struct {
	CurrentPassword string `json:"currentPassword" example:"zaqw@Dekflas;h#&"`
	ConfirmPassword string `json:"confirmPassword" example:"zaqw@Dekflas;h#&"`
	NewPassword     string `json:"newPassword" example:"zaqw@Dekflas;h#&"`
}
type DeleteUser struct {
	Password string `json:"password" example:"zaqw@Dekflas;h#&"`
}
type UpdateUserNotificationsSettings struct {
	DiscordNotificationsSettings bool `json:"discordNotificationsSettings" example:"true"`
	SlackNotificationsSettings   bool `json:"slackNotificationsSettings" example:"true"`
	EmailNotificationsSettings   bool `json:"emailNotificationsSettings" example:"true"`
}
type LoggedUser struct {
	ID      int    `json:"id" example:"11"`
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
}

func NewLoggedUser(id int, email string, name string, surname string) *LoggedUser {
	return &LoggedUser{
		ID:      id,
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
