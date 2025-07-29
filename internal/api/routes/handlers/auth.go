package handlers

import (
	"fmt"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
)

func SetupAuthHadnlers(r *routes.Router) {
	authGroup := r.Group("/auth")
	authGroup.POST("/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hi")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hi from server"))
	})
	authGroup.POST("/login")
	authGroup.GET("/check")
	authGroup.POST("/reset-password")
	authGroup.POST("/reset-password/set-password")
}
