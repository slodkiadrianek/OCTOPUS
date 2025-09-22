package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (a *AppRepository) InsertApp(ctx context.Context, app []DTO.App) error {
	values := ""
	args := make([]any, 0)
	for i := range app {
		values += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		args = append(args, app[i].Id, app[i].Name, app[i].IsDocker, app[i].OwnerID, app[i].IpAddress, app[i].Port)
	}
	values = values[:len(values)-1]
	query := fmt.Sprintf(`INSERT INTO apps (
		id,
		name,
		is_docker,
		owner_id,
		ip_address,
		port
	)  VALUES %s`, values)
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", map[string]any{
			"query": query,
			"args":  app,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add new app to the database")
	}
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		a.Logger.Error("Failed to execute an insert query", map[string]any{
			"query": query,
			"args":  app,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add new app to the database")
	}
	return nil
}

func (a *AppRepository) GetApp(ctx context.Context, id int) (*models.App, error) {
	query := `SELECT * FROM apps WHERE id = $1`
	row := a.Db.QueryRowContext(ctx, query, id)
	var app models.App
	err := row.Scan(&app.Id, &app.Name, &app.IsDocker, &app.OwnerID, &app.SlackWebhook, &app.DiscordWebhook, &app.IpAddress, &app.Port)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"args":  id,
			"err":   err.Error(),
		})
		return nil, err
	}
	return &app, nil
}

// func (a *AppRepository) UpdateApp(ctx context.Context, app DTO.UpdateApp) error {
// 	query := `UPDATE apps SET
// 		name = $2,
// 		description = $3,
// 		db_link = $4,
// 		is_docker = $5,
// 		slack_webhook = $6,
// 		discord_webhook = $7
// 	WHERE id = $1`
// 	stmt, err := a.Db.PrepareContext(ctx, query)
// 	if err != nil {
// 		a.Logger.Error("Failed to prepared statement for execution", map[string]any{
// 			"query": query,
// 			"args":  app,
// 			"err":   err.Error(),
// 		})
// 		return models.NewError(500, "Database", "Failed to update app in the database")
// 	}
// 	_, err = stmt.ExecContext(ctx, app.Id, app.Name, app.Description, app.DbLink, app.ApiLink, app.SlackWebhook, app.DiscordWebhoo)
// 	if err != nil {
// 		a.Logger.Error("Failed to execute an update query", map[string]any{
// 			"query": query,
// 			"args":  app,
// 			"err":   err.Error(),
// 		})
// 	}
// 	return nil
// }

func (a *AppRepository) DeleteApp(ctx context.Context, id int) error {
	query := `DELETE FROM apps WHERE id = $1`
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to delete app from the database")
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		a.Logger.Error("Failed to execute a delete query", map[string]any{
			"query": query,
			"args":  id,
			"err":   err.Error(),
		})
	}
	return nil
}

func (a *AppRepository) GetAppServerAddress(ctx context.Context, id int) (string, error) {
	query := `SELECT apiLink FROM apps WHERE id = $1`
	row := a.Db.QueryRowContext(ctx, query, id)
	var appServerAddress string
	err := row.Scan(&appServerAddress)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"args":  id,
			"err":   err.Error(),
		})
		return "", err
	}
	return appServerAddress, nil
}

func (a *AppRepository) GetDbServerAddress(ctx context.Context, id int) (string, error) {
	query := `SELECT dbLink FROM apps WHERE id = $1`
	row := a.Db.QueryRowContext(ctx, query, id)
	var dbServerAddress string
	err := row.Scan(&dbServerAddress)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"args":  id,
			"err":   err.Error(),
		})
		return "", err
	}
	return dbServerAddress, nil
}

func (a *AppRepository) GetAppStatus(ctx context.Context, id string) (DTO.AppStatus, error) {
	query := "SELECT * FROM apps_statuses WHERE apps_statuses.app_id = $1	"
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepare statement", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return DTO.AppStatus{}, err
	}
	defer stmt.Close()
	var appStatus DTO.AppStatus
	err = stmt.QueryRowContext(ctx, id).Scan(&appStatus.AppId, &appStatus.Status, &appStatus.ChangedAt, &appStatus.Duration)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"args":  id,
			"err":   err.Error(),
		})
		return DTO.AppStatus{}, err
	}
	return appStatus, nil
}

func (a *AppRepository) GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error) {
	query := `SELECT
	    a.id,
	    a.name,
			a.owner_id,
	    a.is_docker,
	    a.ip_address,
	    a.port,
      aps.status
    FROM apps a
		INNER JOIN apps_statuses aps ON a.id = aps.app_id`
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepare statement", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	defer rows.Close()
	apps := make([]*models.AppToCheck, 0)
	for rows.Next() {
		app := &models.AppToCheck{}
		err := rows.Scan(&app.Id, &app.Name, &app.OwnerID, &app.IsDocker, &app.IpAddress, &app.Port, &app.Status)
		if err != nil {
			a.Logger.Error("Failed to scan row", map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return nil, err
		}
		apps = append(apps, app)
	}
	if err := rows.Err(); err != nil {
		a.Logger.Error("Failed to iterate over rows", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	return apps, nil
}

func (a *AppRepository) InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	values := ""
	args := make([]any, 0)
	for i := range appsStatuses {
		values += fmt.Sprintf("($%d,$%d,$%d,$%d),", i*4+1, i*4+2, i*4+3, i*4+4)
		args = append(args, appsStatuses[i].AppId, appsStatuses[i].Status, appsStatuses[i].ChangedAt, appsStatuses[i].Duration)
	}
	values = values[:len(values)-1]

	query := fmt.Sprintf(`
    INSERT INTO apps_statuses(
        app_id,
        status,
        changed_at,
        duration
    ) VALUES %s
    ON CONFLICT (app_id) 
    DO UPDATE SET
        status = EXCLUDED.status,
        changed_at = EXCLUDED.changed_at,
        duration = EXCLUDED.duration
`, values)

	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add app statuses to the database")
	}
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		a.Logger.Error("Failed to execute an insert query", map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add app statuses to the database")
	}
	return nil
}

func (a *AppRepository) GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.SendNotificationTo, error) {
	values := ""
	args := make([]any, 0)
	for i := range appsStatuses {
		values += fmt.Sprintf("$%d,", i+1)
		args = append(args, appsStatuses[i].AppId)
	}
	if len(values) > 0 {
		values = values[:len(values)-1]
	}
	query := fmt.Sprintf(`
	SELECT
    a.id,
    a.name,
		aps.status,
    a.discord_webhook,
    a.slack_webhook,
    u.email,
    u.email_notifications,
    u.discord_notifications,
    u.slack_notifications
	FROM apps a
	INNER JOIN apps_statuses aps ON aps.app_id = a.id
	INNER JOIN users u ON u.id = a.owner_id
	WHERE a.id IN (%s)
	`, values)
	stmt, err := a.Db.PrepareContext(ctx, query)
	if err != nil {
		a.Logger.Error("Failed to prepared statement for execution", map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return []models.SendNotificationTo{}, models.NewError(500, "Database", "Failed to get app from  the database")
	}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		a.Logger.Error("Failed to execute a select query", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	defer rows.Close()
	var dataToSendNotifications []models.SendNotificationTo
	for rows.Next() {
		var objectToSendNotification models.SendNotificationTo
		err := rows.Scan(&objectToSendNotification.Id, &objectToSendNotification.Name, &objectToSendNotification.Status, &objectToSendNotification.DiscordWebhook, &objectToSendNotification.SlackWebhook, &objectToSendNotification.Email, &objectToSendNotification.EmailNotifications, &objectToSendNotification.DiscordNotifications, &objectToSendNotification.SlackNotifications)
		if err != nil {
			a.Logger.Error("Failed to scan row", map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return nil, err
		}
		dataToSendNotifications = append(dataToSendNotifications, objectToSendNotification)
	}
	if err := rows.Err(); err != nil {
		a.Logger.Error("Failed to iterate over rows", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	return dataToSendNotifications, nil
}
