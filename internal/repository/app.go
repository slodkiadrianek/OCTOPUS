package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type AppRepository struct {
	db            *sql.DB
	loggerService utils.LoggerService
}

func NewAppRepository(db *sql.DB, loggerService utils.LoggerService) *AppRepository {
	return &AppRepository{
		db:            db,
		loggerService: loggerService,
	}
}

func (a *AppRepository) InsertApp(ctx context.Context, app []DTO.App) error {
	placeholders := make([]string, 0, len(app))
	args := make([]any, 0, len(app))
	for i := range app {
		preparedValues := fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		args = append(args, app[i].ID, app[i].Name, app[i].IsDocker, app[i].OwnerID, app[i].IpAddress, app[i].Port)
		placeholders = append(placeholders, preparedValues)
	}

	query := fmt.Sprintf(`INSERT INTO apps (
		id,
		name,
		is_docker,
		owner_id,
		ip_address,
		port
	)  VALUES %s`, strings.Join(placeholders, ","))

	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args":  app,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add new app to the database")
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		a.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": query,
			"args":  app,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add new app to the database")
	}

	return nil
}

func (a *AppRepository) GetApp(ctx context.Context, appId string, ownerId int) (*models.App, error) {
	query := `SELECT  
		id,
		name,
		COALESCE(description, ''),
		is_docker,
		owner_id,
		COALESCE(slack_webhook_url, ''),
		COALESCE(discord_webhook_url, ''),
		ip_address,
		port
	FROM apps 
	WHERE id = $1 AND owner_id = $2`
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args": map[string]any{
				"appId":   appId,
				"ownerId": ownerId,
			},
			"err": err.Error(),
		})
		return &models.App{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	var app models.App
	row := stmt.QueryRowContext(ctx, appId, ownerId)
	err = row.Scan(&app.ID, &app.Name, &app.Description, &app.IsDocker, &app.OwnerID, &app.SlackWebhookUrl,
		&app.DiscordWebhookUrl, &app.IpAddress, &app.Port)
	if err != nil {
		a.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"args":  appId,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	return &app, nil
}

func (a *AppRepository) GetApps(ctx context.Context, ownerId int) ([]models.App, error) {
	query := `SELECT 
		id, 
		name, 
		COALESCE(description, ''), 
		is_docker,
		owner_id,
		COALESCE(slack_webhook_url, ''),
		COALESCE(discord_webhook_url, ''),
		ip_address,
		port 
	FROM apps 
	WHERE owner_id = $1`
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args": map[string]any{
				"ownerId": ownerId,
			},
			"err": err.Error(),
		})
		return []models.App{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	rows, err := stmt.QueryContext(ctx, ownerId)
	if err != nil {
		a.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer rows.Close()

	var apps []models.App
	for rows.Next() {
		var app models.App
		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.Description,
			&app.IsDocker,
			&app.OwnerID,
			&app.SlackWebhookUrl,
			&app.DiscordWebhookUrl,
			&app.IpAddress,
			&app.Port,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (a *AppRepository) DeleteApp(ctx context.Context, appId string, ownerId int) error {
	query := `DELETE FROM apps WHERE id = $1 AND owner_id = $2 `
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to delete app from the database")
	}

	_, err = stmt.ExecContext(ctx, appId, ownerId)
	if err != nil {
		a.loggerService.Error(failedToExecuteDeleteQuery, map[string]any{
			"query": query,
			"args":  appId,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to delete app from database")
	}

	return nil
}

func (a *AppRepository) GetAppStatus(ctx context.Context, appId string, ownerId int) (DTO.AppStatus, error) {
	query := "SELECT * FROM apps_statuses WHERE apps_statuses.app_id = $1 AND owner_id = $2	"
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return DTO.AppStatus{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	defer stmt.Close()
	var appStatus DTO.AppStatus
	err = stmt.QueryRowContext(ctx, appId, ownerId).Scan(&appStatus.AppID, &appStatus.Status, &appStatus.ChangedAt,
		&appStatus.Duration)
	if err != nil {
		a.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"args":  appId,
			"err":   err.Error(),
		})
		return DTO.AppStatus{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	return appStatus, nil
}

func (a *AppRepository) GetAppsToCheck(ctx context.Context) ([]*models.AppToCheck, error) {
	query := `
	SELECT
	    a.id,
	    a.name,
		a.owner_id,
	    a.is_docker,
	    a.ip_address,
	    a.port,
		COALESCE(aps.status, 'stopped')
    FROM apps a
		LEFT JOIN apps_statuses aps ON a.id = aps.app_id`
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		a.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer rows.Close()

	apps := make([]*models.AppToCheck, 0)
	for rows.Next() {
		app := &models.AppToCheck{}
		err := rows.Scan(&app.ID, &app.Name, &app.OwnerID, &app.IsDocker, &app.IpAddress, &app.Port, &app.Status)
		if err != nil {
			a.loggerService.Error(failedToScanRows, map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		apps = append(apps, app)
	}

	if err := rows.Err(); err != nil {
		a.loggerService.Error(failedToIterateOverRows, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}

	return apps, nil
}

func (a *AppRepository) UpdateApp(ctx context.Context, appId string, app DTO.UpdateApp, ownerId int) error {
	query := `
	UPDATE apps SET 
        name = $1,
        description = $2,
        ip_address = $3,
        port = $4,
        discord_webhook_url = $5,
		slack_webhook_url = $6 
	WHERE 
	    id = $7 
	  	AND owner_id = $8 `
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update app settings")
	}

	_, err = stmt.ExecContext(ctx, app.Name, app.Description, app.IpAddress, app.Port, app.DiscordWebhookUrl,
		app.SlackWebhookUrl, appId, ownerId)
	if err != nil {
		a.loggerService.Error(failedToExecuteUpdateQuery, map[string]any{
			"query": query,
			"args": map[string]any{
				"appId":   appId,
				"appInfo": app,
			},
			"err": err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update app settings")
	}

	return nil
}

func (a *AppRepository) InsertAppStatuses(ctx context.Context, appsStatuses []DTO.AppStatus) error {
	placeholders := make([]string, 0, len(appsStatuses))
	args := make([]any, 0, len(appsStatuses))
	for i := range appsStatuses {
		preparedValues := fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		args = append(args, appsStatuses[i].AppID, appsStatuses[i].Status, appsStatuses[i].ChangedAt,
			appsStatuses[i].Duration)
		placeholders = append(placeholders, preparedValues)
	}

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
`, strings.Join(placeholders, ","))

	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add app statuses to the database")
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		a.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to add app statuses to the database")
	}

	return nil
}

func (a *AppRepository) GetUsersToSendNotifications(ctx context.Context, appsStatuses []DTO.AppStatus) ([]models.NotificationInfo, error) {
	placeholders := make([]string, 0, len(appsStatuses))
	args := make([]any, 0, len(appsStatuses))
	for i := range appsStatuses {
		preparedValues := fmt.Sprintf("$%d", i+1)
		args = append(args, appsStatuses[i].AppID)
		placeholders = append(placeholders, preparedValues)
	}

	query := fmt.Sprintf(`
	SELECT
		a.id,
		a.name,
		aps.status,
		COALESCE(a.discord_webhook_url,''),
		COALESCE(a.slack_webhook_url,''),
		u.email,
		u.email_notifications_settings,
		u.discord_notifications_settings,
		u.slack_notifications_settings
	FROM apps a
		INNER JOIN apps_statuses aps ON aps.app_id = a.id
		INNER JOIN users u ON u.id = a.owner_id
	WHERE a.id IN (%s)
	AND (
            (u.discord_notifications_settings = true AND a.discord_webhook_url  != '')
            OR (u.email_notifications_settings = true )
            OR (u.slack_notifications_settings = true AND a.slack_webhook_url != '')
        )
	`, strings.Join(placeholders, ","))
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		a.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args":  appsStatuses,
			"err":   err.Error(),
		})
		return []models.NotificationInfo{}, models.NewError(500, "Database", "Failed to get app from the database")
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		a.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, err
	}
	defer rows.Close()

	var notifications []models.NotificationInfo
	for rows.Next() {
		var notification models.NotificationInfo
		err := rows.Scan(&notification.ID, &notification.Name, &notification.Status, &notification.DiscordWebhookUrl,
			&notification.SlackWebhookUrl, &notification.Email, &notification.EmailNotificationsSettings,
			&notification.DiscordNotificationsSettings, &notification.SlackNotificationsSettings)
		if err != nil {
			a.loggerService.Error(failedToScanRows, map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return nil, models.NewError(500, "Database", "Failed to get app from the database")
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		a.loggerService.Error(failedToIterateOverRows, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return nil, models.NewError(500, "Database", "Failed to get app from the database")
	}

	return notifications, nil
}
