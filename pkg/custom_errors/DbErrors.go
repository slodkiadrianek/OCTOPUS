package custom_errors

const (
	ERR_DB_NOT_FOUND = "NOT_FOUND"
	ERR_DB_SERVER    = "DB_FAILED"
)

var Err_db_not_found = map[string]string{
	"errorCategory":    ERR_DB_NOT_FOUND,
	"errorDescription": "Failed to find any rows",
}

var Err_db_server_insert = map[string]string{
	"errorCategory":    ERR_DB_SERVER,
	"errorDescription": "Problems occured during inserting data to database",
}

var Err_db_server_select = map[string]string{
	"errorCategory":    ERR_DB_SERVER,
	"errorDescription": "Problems occured during selecting data from database",
}

var Err_db_server_update = map[string]string{
	"errorCategory":    ERR_DB_SERVER,
	"errorDescription": "Problems occured during updating data from database",
}

var Err_db_server_Delete = map[string]string{
	"errorCategory":    ERR_DB_SERVER,
	"errorDescription": "Problems occured during deleting data from database",
}
