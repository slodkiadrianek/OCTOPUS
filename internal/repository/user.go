package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type UserRepository struct {
	Db            *sql.DB
	LoggerService *utils.Logger
}

func NewUserRepository(db *sql.DB, loggerService *utils.Logger) *UserRepository {
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
	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.Name, &user.Surname, &user.Password, &user.DiscordNotifications, &user.EmailNotifications, &user.SlackNotifications, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.LoggerService.Info("user not found", map[string]any{
				"email": email,
			})
			return models.User{
				Id: 0,
			}, nil
		}
		u.LoggerService.Error("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{email},
			"error": err,
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
		u.LoggerService.Info("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{user.Name, user.Surname, user.Email, password},
			"err":   err.Error(),
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
		return models.NewError(500, "Database", "Failed to update data in database")
	}
	_, err = stmt.ExecContext(ctx, user.Name, user.Surname, user.Email, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{user.Name, user.Surname, user.Email, userId},
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update data in database")
	}
	return nil
}

func (u *UserRepository) UpdateUserNotifications(ctx context.Context, userId int, userNotifications DTO.UpdateUserNotifications,
) error {
	query := `UPDATE users SET discord_notifications=$1, slack_notifications=$2, email_notifications=$3 WHERE id=$4`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", query)
		return models.NewError(500, "Database", "Failed to update data in database")
	}

	_, err = stmt.ExecContext(ctx, userNotifications.DiscordNotifications, userNotifications.SlackNotifications, userNotifications.EmailNotifications, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{userNotifications, userId},
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update data in database")
	}
	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, password string, userId int) error {
	query := `DELETE FROM users WHERE id=$1`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to delete data from database")
	}

	_, err = stmt.ExecContext(ctx, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{password, userId},
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to delete data from database")
	}
	return nil
}

func (u *UserRepository) FindUserById(ctx context.Context, userId int) (models.User, error) {
	query := `
	SELECT * FROM users WHERE id = $1`
	u.Db.SetMaxOpenConns(1000)
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.User{}, models.NewError(500, "Database", "Failed to get data from the database")
	}
	defer stmt.Close()
	var user models.User
	err = stmt.QueryRowContext(ctx, userId).Scan(&user.Id, &user.Email, &user.Name, &user.Surname, &user.Password,
		&user.DiscordNotifications, &user.EmailNotifications, &user.SlackNotifications, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.LoggerService.Info("user not found", map[string]any{
				"userId": userId,
			})
			return models.User{
				Id: 0,
			}, nil
		}
		u.LoggerService.Error("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{userId},
			"err":   err.Error(),
		})
		return models.User{}, models.NewError(500, "Database", "Failed to get data from the database")
	}
	return user, nil
}

func (u *UserRepository) ChangeUserPassword(ctx context.Context, userId int, newPassword string) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`
	stmt, err := u.Db.PrepareContext(ctx, query)
	if err != nil {
		u.LoggerService.Info("failed to prepare query for execution", query)
		return models.NewError(500, "Database", "Failed to update data in database")
	}
	_, err = stmt.ExecContext(ctx, newPassword, userId)
	if err != nil {
		u.LoggerService.Info("failed to execute query for execution", map[string]any{
			"query": query,
			"args":  []any{newPassword, userId},
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update data in database")
	}
	return nil
}
