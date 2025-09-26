package repository

import (
	"database/sql"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteRepository struct {
	Db     *sql.DB
	Logger *utils.Logger
}

func NewRouteRepository(db *sql.DB, logger *utils.Logger) *RouteRepository {
	return &RouteRepository{
		Db:     db,
		Logger: logger,
	}
}

// func (r *RouteRepository) InsertRoutes(ctx context.Context, routes []DTO.CreateRoute) error {
// 	insertQuery := `INSERT INTO routes (
// 		method,
// 		path,
// 		queryData,
// 		paramData,
// 		bodyData,
// 		expectedBodyData,
// 		expectedStatusCode
// 	) VALUES`
// 	placeholders := []string{}
// 	for i, _ := range routes {
// 		values := fmt.Sprintf(" (%s,%s,%s,%s,%s,%s,%s)", "$"+string(i*10), "$"+string(i*10+1), "$"+string(i*10+2), "$"+string(i*10+3), "$"+string(i*10+4), "$"+string(i*10+5), "$"+string(i*10+6))
// 		placeholders = append(placeholders, values)
// 	}
// 	insertQuery += strings.Join(placeholders, ", ")

// 	return r.Db.ExecContext(ctx, insertQuery).Error
// }
