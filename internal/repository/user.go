package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/schema"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type UserRepository struct {
	Db            *sql.DB
	LoggerService *logger.Logger
}

func NewUserRepository(db *sql.DB, loggerService *logger.Logger) *UserRepository {
	return &UserRepository{
		Db:            db,
		LoggerService: loggerService,
	}
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", err)
		return models.User{}, models.NewError(500, "Database", "Failed to get data from the database")
	}
	defer stmt.Close()
	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Name, &user.Surname, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.LoggerService.Info("user not found", map[string]interface{}{
				"email": email,
			})
			return models.User{
				Id: 0,
			}, nil
		}
		u.LoggerService.Info("failed to execute query for execution", map[string]interface{}{
			"query": query,
			"args":  []interface{}{email},
		})
		return models.User{}, models.NewError(500, "Database", "Failed to get data from the database")
	}
	return user, nil
}

func (u *UserRepository) InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error {
	query := `INSERT INTO users(name, surname, email, password) VALUES($1, $2, $3, $4 )`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", err)
		return models.NewError(500, "Database", "Failed to insert data to the database")
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, user.Name, user.Surname, user.Email, password)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]interface{}{
			"query": query,
			"args":  []interface{}{user.Name, user.Surname, user.Email, password},
		})
		return models.NewError(500, "Database", "Failed to insert data to the database")
	}

	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error {
	query := `UPDATE users SET name=$1, surname=$2, email=$3 WHERE id=$4`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", query)
		return models.NewError(500, "Database", "Failed to update data in  database")
	}
	_, err = stmt.ExecContext(ctx, user.Name, user.Surname, user.Email, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]interface{}{
			"query": query,
			"args":  []interface{}{user.Name, user.Surname, user.Email, userId},
		})
		return models.NewError(500, "Database", "Failed to insert data to the database")
	}
	return nil
}

func (u *UserRepository) UpdateUserNotifications(ctx context.Context, userId int, userNotifications schema.UpdateUserNotifications) error {
	query := `UPDATE users SET discordNotifications=$1, slackNotifications=$2, emailNotifications=$3 WHERE id=$4`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", query)
		return models.NewError(500, "Database", "Failed to update data in  database")
	}

	_, err = stmt.ExecContext(ctx, userNotifications.DiscordNotifications, userNotifications.SlackNotifications, userNotifications.EmailNotifications, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]interface{}{
			"query": query,
			"args":  []interface{}{userNotifications, userId},
		})
		return models.NewError(500, "Database", "Failed to insert data to the database")
	}
	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, password string, userId int) error {
	query := `DELETE FROM users WHERE id=$1 AND password=$2`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", query)
		return models.NewError(500, "Database", "Failed to update data in  database")
	}

	_, err = stmt.ExecContext(ctx, userId, password)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]interface{}{
			"query": query,
			"args":  []interface{}{password, userId},
		})
		return models.NewError(500, "Database", "Failed to insert data to the database")
	}
	return nil
}
