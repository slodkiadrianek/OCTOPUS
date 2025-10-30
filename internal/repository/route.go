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
	Db            *sql.DB
	LoggerService *utils.Logger
}

func NewRouteRepository(db *sql.DB, logger *utils.Logger) *RouteRepository {
	return &RouteRepository{
		Db:            db,
		LoggerService: logger,
	}
}

func (r *RouteRepository) UpdateWorkingRoutesStatuses(ctx context.Context, routesStatuses map[int]string) error {
	placeholders := []string{}
	argPos := 1
	args := make([]any, 0)
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
	_, err := r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		r.LoggerService.Error("Failed to to update routes statuses", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return models.NewError(500, "Database", "Failed to update routes statuses")
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
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return []models.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.LoggerService.Error("Failed to execute query", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return []models.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var routesToTest []models.RouteToTest
	for rows.Next() {
		var routeToTest models.RouteToTest
		err := rows.Scan(&routeToTest.ID, &routeToTest.IpAddress, &routeToTest.Port, &routeToTest.Name,
			&routeToTest.AppId,
			&routeToTest.ParentID, &routeToTest.Status,
			&routeToTest.Path,
			&routeToTest.Method, &routeToTest.RequestAuthorization, &routeToTest.RequestQuery, &routeToTest.RequestParams, &routeToTest.RequestBody, &routeToTest.NextRouteBody, &routeToTest.NextRouteParams, &routeToTest.NextRouteQuery, &routeToTest.NextAuthorizationHeader, &routeToTest.ResponseStatusCode, &routeToTest.ResponseBody)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return []models.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		routesToTest = append(routesToTest, routeToTest)
	}
	return routesToTest, nil
}

func (r *RouteRepository) InsertRoutesInfo(ctx context.Context, routesInfo []*DTO.RouteInfo) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
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
	stmt, err := r.Db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var routesInfoIds []int
	for rows.Next() {
		var routeInfoId int
		err := rows.Scan(&routeInfoId)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		routesInfoIds = append(routesInfoIds, routeInfoId)
	}
	return routesInfoIds, nil
}

func (r *RouteRepository) InsertRoutesRequests(ctx context.Context,
	routesRequests []*DTO.RouteRequest,
) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
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
	stmt, err := r.Db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var routesRequestsIds []int
	for rows.Next() {
		var routeRequestId int
		err := rows.Scan(&routeRequestId)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		routesRequestsIds = append(routesRequestsIds, routeRequestId)
	}
	return routesRequestsIds, nil
}

func (r *RouteRepository) InsertRoutesResponses(ctx context.Context,
	routesResponses []*DTO.RouteResponse) ([]int,
	error,
) {
	placeholders := []string{}
	args := make([]any, 0)
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
	stmt, err := r.Db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var routesResponsesIds []int
	for rows.Next() {
		var routeResponseId int
		err := rows.Scan(&routeResponseId)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		routesResponsesIds = append(routesResponsesIds, routeResponseId)
	}
	return routesResponsesIds, nil
}

func (r *RouteRepository) InsertNextRoutesData(ctx context.Context,
	nextRoutes []*DTO.NextRoute) ([]int,
	error,
) {
	placeholders := []string{}
	args := make([]any, 0)
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
	stmt, err := r.Db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
		return []int{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var NextRoutesDataIds []int
	for rows.Next() {
		var nextRouteDataId int
		err := rows.Scan(&nextRouteDataId)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": insertQuery,
				"err":   err.Error(),
			})
			return []int{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		NextRoutesDataIds = append(NextRoutesDataIds, nextRouteDataId)
	}
	return NextRoutesDataIds, nil
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
	stmt, err := r.Db.PrepareContext(ctx, insertQuery)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
		})
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRowContext(ctx, workingRoute.Name, workingRoute.AppId, workingRoute.ParentID, workingRoute.RouteID,
		workingRoute.RequestID, workingRoute.ResponseID, workingRoute.NextRouteDataId, workingRoute.Status).Scan(&id)
	if err != nil {
		r.LoggerService.Error("Failed to prepare statement", map[string]any{
			"query": insertQuery,
			"err":   err.Error(),
			"data":  workingRoute,
		})
		return 0, models.NewError(500, "Database", "Failed to get data from database")
	}
	return id, nil
}
