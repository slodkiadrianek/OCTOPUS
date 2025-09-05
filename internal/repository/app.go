package repository

import (
	"context"
	"database/sql"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils/logger"
)

type AppRepository struct {
	Db     *sql.DB
	Logger *logger.Logger
}

func NewAppRepository(db *sql.DB, logger *logger.Logger) *AppRepository {
	return &AppRepository{
		Db:     db,
		Logger: logger,
	}
}

func (a *AppRepository) InsertApp(ctx context.Context, app DTO.App) error {
	query := `INSERT INTO apps(
		name,
		description,
		dbLink,
		apiLink,
		ownerId,
		slackWebhook,
		discordWebhook
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)`
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", query)
		return models.NewError(500, "Database", "Failed to add new app to the database")
	}
	_, err = stmt.ExecContext(ctx, app.Name, app.DbLink, app.ApiLink, app.OwnerID, app.SlackWebhook, app.DiscordWebhook)
	if err != nil {
		a.Logger.Error("Failed to execute an insert query", map[string]interface{}{
			"query": query,
			"args":  app,
		})
	}
	return nil
}

func (a *AppRepository) GetApp(ctx context.Context, id int) (*models.App, error) {
	query := `SELECT * FROM apps WHERE id = $1`
	row := a.Db.QueryRowContext(ctx, query, id)
	var app models.App
	err := row.Scan(&app.Id, &app.Name, &app.DbLink, &app.ApiLink, &app.OwnerID, &app.SlackWebhook, &app.DiscordWebhook)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"args":  id,
		})
		return nil, err
	}
	return &app, nil
}

func (a *AppRepository) UpdateApp(ctx context.Context, app DTO.UpdateApp) error {
	query := `UPDATE apps SET
		name = $2,
		description = $3,
		dbLink = $4,
		apiLink = $5,
		slackWebhook = $6,
		discordWebhook = $7
	WHERE id = $1`
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", query)
		return models.NewError(500, "Database", "Failed to update app in the database")
	}
	_, err = stmt.ExecContext(ctx, app.Id, app.Name, app.Description, app.DbLink, app.ApiLink, app.SlackWebhook, app.DiscordWebhook)
	if err != nil {
		a.Logger.Error("Failed to execute an update query", map[string]any{
			"query": query,
			"args":  app,
		})
	}
	return nil
}

func (a *AppRepository) DeleteApp(ctx context.Context, id int) error {
	query := `DELETE FROM apps WHERE id = $1`
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", query)
		return models.NewError(500, "Database", "Failed to delete app from the database")
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		a.Logger.Error("Failed to execute a delete query", map[string]any{
			"query": query,
			"args":  id,
		})
	}
	return nil
}
