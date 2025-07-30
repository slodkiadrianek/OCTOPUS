package handlers

import (
	"fmt"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/api/routes"
)

func SetupAuthHadnlers(r *routes.Router) {
	authGroup := r.Group("/api/v1/auth")
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

func SetupWorkspaceHandlers(r *routes.Router) {
	workspaceGroup := r.Group("/api/v1/workspaces")
	workspaceGroup.POST("")
	workspaceGroup.DELETE("/:id")
}
