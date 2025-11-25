package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type userClaims struct {
	Id      int    `json:"id" example:"11"`
	Email   string `json:"email" example:"joedoe@email.com"`
	Name    string `json:"name" example:"Joe"`
	Surname string `json:"surname" example:"Doe"`
	exp     int64
	jwt.RegisteredClaims
}

type JWT struct {
	token         string
	loggerService utils.LoggerService
	cacheService  interfaces.CacheService
}

func NewJWT(token string, loggerService utils.LoggerService, cacheService interfaces.CacheService) *JWT {
	return &JWT{
		token:         token,
		loggerService: loggerService,
		cacheService:  cacheService,
	}
}

func (j JWT) GenerateToken(user DTO.LoggedUser) (string, error) {
	j.loggerService.Info("started signing a new token")
	tokenWithData := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	})

	tokenString, err := tokenWithData.SignedString([]byte(j.token))
	if err != nil {
		j.loggerService.Error("Failed to sign token properly", err)
		return "", models.NewError(401, "Authorization", "Failed to login the user")
	}

	j.loggerService.Info("Successfully signed a new token")
	return tokenString, nil
}

func (j JWT) parseClaimsFromToken(user *userClaims, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, user, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.token), nil
	})
	return token, err
}

func (j JWT) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			j.loggerService.Info("token is missing", authHeader)
			err := models.NewError(401, "Authorization", "Failed to authorize a user")
			utils.SetError(w, r, err)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		result, err := j.cacheService.ExistsData(r.Context(), "blacklist-"+tokenString)
		if err != nil {
			j.loggerService.Info("Failed to check blacklist", err)
			err := models.NewError(401, "Authorization", "Failed to check blacklist")
			utils.SetError(w, r, err)
			return

		}

		if result > 0 {
			j.loggerService.Info("Token is blacklisted", tokenString)
			err := models.NewError(401, "Authorization", "Token is blacklisted")
			utils.SetError(w, r, err)
			return

		}

		var user *userClaims
		tokenWithData, err := j.parseClaimsFromToken(user, tokenString)

		if err != nil {
			j.loggerService.Info("Failed to read data properly", err)
			err := models.NewError(401, "Authorization", "Provided token is invalid")
			utils.SetError(w, r, err)
			return

		}

		if !tokenWithData.Valid {
			j.loggerService.Info("Provided token is invalid", tokenString)
			err := models.NewError(401, "Authorization", "Provided token is invalid")
			utils.SetError(w, r, err)
			return

		}

		r = utils.SetContext(r, "id", user.Id)

		r = utils.SetContext(r, "email", user.Email)
		next.ServeHTTP(w, r)
	})
}

func (j JWT) BlacklistUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			j.loggerService.Info("token is missing", authHeader)
			err := models.NewError(401, "Authorization", "Failed to authorize a user")
			utils.SetError(w, r, err)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		result, err := j.cacheService.ExistsData(r.Context(), "blacklist-"+tokenString)
		if err != nil {
			j.loggerService.Info("Failed to check blacklist", err)
			err := models.NewError(401, "Authorization", "Failed to check blacklist")
			utils.SetError(w, r, err)
			return
		}

		if result > 0 {
			j.loggerService.Info("Token is blacklisted", tokenString)
			err := models.NewError(401, "Authorization", "Token is blacklisted")
			utils.SetError(w, r, err)
			return
		}

		var user *userClaims
		tokenWithData, err := j.parseClaimsFromToken(user, tokenString)

		if err != nil {
			j.loggerService.Info("Failed to read data properly", tokenString)
			err := models.NewError(401, "Authorization", "Failed to read token")
			utils.SetError(w, r, err)
			return
		}

		if !tokenWithData.Valid {
			j.loggerService.Info("Provided token is invalid", tokenString)
			err := models.NewError(401, "Authorization", "Provided token is invalid")
			utils.SetError(w, r, err)
			return
		}

		expirationTime := time.Until(user.ExpiresAt.Time)
		err = j.cacheService.SetData(r.Context(), "blacklist-"+tokenString, "true", expirationTime)
		if err != nil {
			j.loggerService.Info("Failed to set data in cache", err)
			err := models.NewError(401, "Authorization", "Failed to authorize a user")
			utils.SetError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
