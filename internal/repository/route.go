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

func (r *RouteRepository) GetWorkingRoutesToTest(ctx context.Context) ([]DTO.RouteToTest, error) {
	query := `
SELECT
		wr.id,
    a.ip_address,
    a.port,
    wr.app_id,
    wr.parent_id,
    wr.status,
    r.path,
    r.method,
    rr.authorization_header,
    rr.query_data,
    rr.param_data,
    rr.body_data,
    nrd.next_route_body,
    nrd.next_route_params,
    nrd.next_route_query,
		nrd.next_authorization_header,
    re.status_code,
    re.body_data
FROM working_routes wr
    INNER JOIN public.routes r on wr.route_id = r.id
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
		return []DTO.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		r.LoggerService.Error("Failed to execute query", map[string]any{
			"query": query,
			"err":   err.Error(),
		})
		return []DTO.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
	}
	var routesToTest []DTO.RouteToTest
	for rows.Next() {
		var routeToTest DTO.RouteToTest
		err := rows.Scan(&routeToTest.Id, &routeToTest.IpAddress, &routeToTest.Port, &routeToTest.Name, &routeToTest.AppId,
			&routeToTest.ParentId, &routeToTest.Status,
			&routeToTest.Path,
			&routeToTest.Method, &routeToTest.RequestQuery, &routeToTest.RequestParams, &routeToTest.RequestBody, &routeToTest.NextRouteBody, &routeToTest.NextRouteParams, &routeToTest.NextRouteQuery, &routeToTest.NextAuthorizationHeader, &routeToTest.ResponseStatusCode, &routeToTest.ResponseBody)
		if err != nil {
			r.LoggerService.Error("Failed to scan row", map[string]any{
				"query": query,
				"err":   err.Error(),
			})
			return []DTO.RouteToTest{}, models.NewError(500, "Database", "Failed to get data from database")
		}
		routesToTest = append(routesToTest, routeToTest)
	}
	return routesToTest, nil
}

func (r *RouteRepository) InsertRoutesInfo(ctx context.Context, routesInfo []*DTO.RouteInfo) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
	for i := range routesInfo {
		values := fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2)
		placeholders = append(placeholders, values)
		args = append(args, routesInfo[i].Path, routesInfo[i].Method)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO routes (
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

func (r *RouteRepository) InsertRoutesReuquests(ctx context.Context, routesRequests []*DTO.RouteRequest) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
	for i := range routesRequests {
		values := fmt.Sprintf("($%d::jsonb,$%d::jsonb,$%d::jsonb,$%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		placeholders = append(placeholders, values)
		args = append(args, routesRequests[i].RequestBody, routesRequests[i].RequestParams, routesRequests[i].RequestQuery, routesRequests[i].RequestAuthorization)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO routes_requests (
	body_data,
	param_data,
	query_data,
	authorization_header
	) VALUES
	%s
	ON CONFLICT(	
		body_data,
		param_data,
		query_data,
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

func (r *RouteRepository) InsertRoutesResponses(ctx context.Context, routesResponses []*DTO.RouteResponse) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
	for i := range routesResponses {
		values := fmt.Sprintf("($%d,$%d::jsonb)", i*2+1, i*2+2)
		placeholders = append(placeholders, values)
		args = append(args, routesResponses[i].ResponseStatusCode, routesResponses[i].ResponseBody)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO routes_responses (
	status_code,
	body_data
	) VALUES
	%s
	ON CONFLICT(	
		status_code,
		body_data
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

func (r *RouteRepository) InsertNextRoutesData(ctx context.Context, nextRoutesData []*DTO.NextRouteData) ([]int, error) {
	placeholders := []string{}
	args := make([]any, 0)
	for i := range nextRoutesData {
		values := fmt.Sprintf("($%d::jsonb,$%d::jsonb,$%d::jsonb, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		placeholders = append(placeholders, values)
		args = append(args, nextRoutesData[i].NextRouteBody, nextRoutesData[i].NextRouteParams, nextRoutesData[i].NextRouteQuery, nextRoutesData[i].NextAuthorizationHeader)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO next_route_data (
	next_route_body,
	next_route_params,
	next_route_query,
	next_route_authorization_header
	) VALUES
	%s
	ON CONFLICT(	
	next_route_body,
	next_route_params,
	next_route_query,
	next_route_authorization_header
	)
	DO UPDATE
		SET next_route_query = EXCLUDED.next_route_query
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

func (r *RouteRepository) InsertWorkingRoute(ctx context.Context, workingRoute DTO.WorkingRoute) (int, error) {
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
	err = stmt.QueryRowContext(ctx, workingRoute.Name, workingRoute.AppId, workingRoute.ParentId, workingRoute.RouteId,
		workingRoute.RequestId, workingRoute.ResponseId, workingRoute.NextRouteDataId, workingRoute.Status).Scan(&id)
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
