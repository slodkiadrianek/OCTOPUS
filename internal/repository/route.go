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

type RouteRepository struct {
	db            *sql.DB
	loggerService utils.LoggerService
}

func NewRouteRepository(db *sql.DB, loggerService utils.LoggerService) *RouteRepository {
	return &RouteRepository{
		db:            db,
		loggerService: loggerService,
	}
}

func (r *RouteRepository) CheckRouteStatus(ctx context.Context, routeID int) (string, error) {
	query := "SELECT status FROM working_routes WHERE id = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args":  routeID,
			"err":   err.Error(),
		})
		return "", models.NewError(500, "Database", "failed to check route status")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()
	var routeStatus string
	err = stmt.QueryRowContext(ctx, routeID).Scan(&routeStatus)
	if err != nil {
		r.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"args":  routeID,
			"err":   err.Error(),
		})
		return "", models.NewError(500, "Database", "failed to check route status")
	}
	return routeStatus, nil
}

func (r *RouteRepository) UpdateWorkingRoutesStatuses(ctx context.Context, routesStatuses map[int]string) error {
	placeholders := make([]string, 0, len(routesStatuses))
	argPos := 1
	args := make([]any, 0, len(routesStatuses))
	for i, val := range routesStatuses {
		preparedValues := fmt.Sprintf("($%d,$%d)", argPos, argPos+1)
		args = append(args, int(i), val)
		placeholders = append(placeholders, preparedValues)
		argPos += 2
	}

	query := fmt.Sprintf(`
	UPDATE working_routes AS t
	SET 
    status = v.status
	FROM (VALUES
	%s
	) AS v(id, status)
		WHERE t.id = v.id::integer;
	`, strings.Join(placeholders, ","))
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"args":  routesStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "failed to update routes statuses")
	}
	defer func() {
		if err != nil {
			if closeErr := stmt.Close(); closeErr != nil {
				r.loggerService.Error(failedToCloseStatement, closeErr)
			}
		}
	}()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		r.loggerService.Error(failedToExecuteUpdateQuery, map[string]any{
			"query": query,
			"args":  routesStatuses,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "failed to update routes statuses")
	}

	return nil
}

func (r *RouteRepository) GetWorkingRoutesToTest(ctx context.Context) ([]models.RouteToTest, error) {
	query := `
SELECT
	wr.id,
    a.ip_address,
    a.port,
	wr.name,
    wr.app_id,
    wr.parent_id,
    wr.status,
    rf.path,
    rf.method,
    rr.authorization_header,
    rr.query,
    rr.params,
    rr.body,
    nrd.body,
    nrd.params,
    nrd.query,
	nrd.authorization_header,
    re.status_code,
    re.body
FROM working_routes wr
    INNER JOIN public.routes_info rf on wr.route_id = rf.id
    INNER JOIN public.routes_requests rr on wr.request_id = rr.id
    INNER JOIN public.next_route_data nrd on wr.next_route_data_id = nrd.id
    INNER JOIN public.routes_responses re on re.id = wr.response_id
    inner join public.apps a on a.id = wr.app_id
    INNER JOIN apps_statuses aps on aps.app_id = wr.app_id
WHERE aps.status = 'running'
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return []models.RouteToTest{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.loggerService.Error(failedToExecuteSelectQuery, map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return []models.RouteToTest{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.loggerService.Error(failedToScanRow, closeErr)
		}
	}()

	var routesToTest []models.RouteToTest
	for rows.Next() {
		var routeToTest models.RouteToTest
		err := rows.Scan(&routeToTest.ID, &routeToTest.IPAddress, &routeToTest.Port, &routeToTest.Name,
			&routeToTest.AppID,
			&routeToTest.ParentID, &routeToTest.Status,
			&routeToTest.Path,
			&routeToTest.Method, &routeToTest.RequestAuthorization, &routeToTest.RequestQuery, &routeToTest.RequestParams, &routeToTest.RequestBody, &routeToTest.NextRouteBody, &routeToTest.NextRouteParams, &routeToTest.NextRouteQuery, &routeToTest.NextAuthorizationHeader, &routeToTest.ResponseStatusCode, &routeToTest.ResponseBody)
		if err != nil {
			r.loggerService.Error(failedToScanRow, map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return []models.RouteToTest{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		routesToTest = append(routesToTest, routeToTest)
	}

	return routesToTest, nil
}

func (r *RouteRepository) InsertRoutesInfo(ctx context.Context, routesInfo []*DTO.RouteInfo) ([]int, error) {
	placeholders := make([]string, 0, len(routesInfo))
	args := make([]any, 0, len(routesInfo))
	for i := range routesInfo {
		preparedValues := fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2)
		placeholders = append(placeholders, preparedValues)
		args = append(args, routesInfo[i].Path, routesInfo[i].Method)
	}

	insertQuery := fmt.Sprintf(`INSERT INTO routes_info (
		path,
		method
	) VALUES
	%s
	ON CONFLICT(path,method)
	DO UPDATE
		SET path = EXCLUDED.path
	Returning id`, strings.Join(placeholders, ","))

	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseRows, closeErr)
		}
	}()

	routesInfoIDs := make([]int, 0, len(routesInfo))
	for rows.Next() {
		var routeInfoID int
		err := rows.Scan(&routeInfoID)
		if err != nil {
			r.loggerService.Error(failedToScanRow, map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		routesInfoIDs = append(routesInfoIDs, routeInfoID)
	}
	return routesInfoIDs, nil
}

func (r *RouteRepository) InsertRoutesRequests(ctx context.Context,
	routesRequests []*DTO.RouteRequest,
) ([]int, error) {
	placeholders := make([]string, 0, len(routesRequests))
	args := make([]any, 0, len(routesRequests))
	for i := range routesRequests {
		preparedValues := fmt.Sprintf("($%d::jsonb,$%d::jsonb,$%d::jsonb,$%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		placeholders = append(placeholders, preparedValues)
		args = append(args, routesRequests[i].Body, routesRequests[i].Params, routesRequests[i].Query,
			routesRequests[i].AuthorizationHeader)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO routes_requests (
		body,
		params,
		query,
		authorization_header
	) VALUES
	%s
	ON CONFLICT(	
		body,
		params,
		query,
		authorization_header
	)
	DO UPDATE
		SET authorization_header = EXCLUDED.authorization_header
	Returning id`, strings.Join(placeholders, ","))
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseRows, closeErr)
		}
	}()

	routesRequestsIDs := make([]int, 0, len(routesRequests))
	for rows.Next() {
		var routeRequestID int
		err := rows.Scan(&routeRequestID)
		if err != nil {
			r.loggerService.Error(failedToScanRow, map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		routesRequestsIDs = append(routesRequestsIDs, routeRequestID)
	}
	return routesRequestsIDs, nil
}

func (r *RouteRepository) InsertRoutesResponses(ctx context.Context,
	routesResponses []*DTO.RouteResponse) ([]int,
	error,
) {
	placeholders := make([]string, 0, len(routesResponses))
	args := make([]any, 0, len(routesResponses))
	for i := range routesResponses {
		values := fmt.Sprintf("($%d,$%d::jsonb)", i*2+1, i*2+2)
		placeholders = append(placeholders, values)
		args = append(args, routesResponses[i].StatusCode, routesResponses[i].Body)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO routes_responses (
		status_code,
		body
	) VALUES
	%s
	ON CONFLICT(	
		status_code,
		body
	)
	DO UPDATE
		SET status_code = EXCLUDED.status_code
	Returning id`, strings.Join(placeholders, ","))
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseRows, closeErr)
		}
	}()

	routesResponsesIDs := make([]int, 0, len(routesResponses))
	for rows.Next() {
		var routeResponseID int
		err := rows.Scan(&routeResponseID)
		if err != nil {
			r.loggerService.Error(failedToScanRow, map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		routesResponsesIDs = append(routesResponsesIDs, routeResponseID)
	}
	return routesResponsesIDs, nil
}

func (r *RouteRepository) InsertNextRoutesData(ctx context.Context,
	nextRoutes []*DTO.NextRoute) ([]int,
	error,
) {
	placeholders := make([]string, 0, len(nextRoutes))
	args := make([]any, 0, len(nextRoutes))
	for i := range nextRoutes {
		values := fmt.Sprintf("($%d::jsonb,$%d::jsonb,$%d::jsonb, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		placeholders = append(placeholders, values)
		args = append(args, nextRoutes[i].Body, nextRoutes[i].Params, nextRoutes[i].Query, nextRoutes[i].AuthorizationHeader)
	}
	insertQuery := fmt.Sprintf(`
	WITH input_data AS (
    SELECT * FROM (VALUES 
		%s
 ) AS v(body,params,query,authorization_header)
),
upserted AS (
    INSERT INTO next_route_data (
        body, params, query, authorization_header
    )
    SELECT DISTINCT ON (body, params,query,authorization_header)
        body,params,query,authorization_header
    FROM input_data
    ON CONFLICT(body,params,query,authorization_header) 
    DO UPDATE SET body = EXCLUDED.body
    RETURNING *
)
SELECT u.id
FROM input_data i
JOIN upserted u USING (body,params,query,authorization_header);
	`, strings.Join(placeholders, ","))
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseRows, closeErr)
		}
	}()

	nextRoutesDataIDs := make([]int, 0, len(nextRoutes))
	for rows.Next() {
		var nextRouteDataID int
		err := rows.Scan(&nextRouteDataID)
		if err != nil {
			r.loggerService.Error(failedToScanRow, map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", failedToGetDataFromDatabase)
		}
		nextRoutesDataIDs = append(nextRoutesDataIDs, nextRouteDataID)
	}
	return nextRoutesDataIDs, nil
}

func (r *RouteRepository) InsertWorkingRoute(ctx context.Context, workingRoute DTO.WorkingRoute) (int,
	error,
) {
	insertQuery := `INSERT INTO working_routes (
	name,
    app_id,
    parent_id,
    route_id,
    request_id,
    response_id,
    next_route_data_id,
    status
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
ON CONFLICT (
	name,
    app_id,
    parent_id,
    route_id,
    request_id,
    response_id,
    next_route_data_id
)
DO UPDATE SET 
    status = EXCLUDED.status
RETURNING id`
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.loggerService.Error(failedToPrepareQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			r.loggerService.Error(failedToCloseStatement, closeErr)
		}
	}()

	var id int
	err = stmt.QueryRowContext(ctx, workingRoute.Name, workingRoute.AppID, workingRoute.ParentID, workingRoute.RouteID,
		workingRoute.RequestID, workingRoute.ResponseID, workingRoute.NextRouteDataID, workingRoute.Status).Scan(&id)
	if err != nil {
		r.loggerService.Error(failedToExecuteInsertQuery, map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
			"data":  workingRoute,
		})
		return 0, models.NewError(500, "Database", failedToGetDataFromDatabase)
	}
	return id, nil
}
