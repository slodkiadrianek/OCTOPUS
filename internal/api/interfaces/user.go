package interfaces

import "net/http"

type AuthController interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
	VerifyUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
}

type UserController interface {
	GetUserInfo(w http.ResponseWriter, r *http.Request)
	InsertUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	UpdateUserNotifications(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ChangeUserPassword(w http.ResponseWriter, r *http.Request)
}
