package repository

const failedToPrepareQuery = "failed to prepared statement for execution"

const (
	failedToExecuteInsertQuery = "failed to execute an insert query"
	failedToExecuteSelectQuery = "failed to execute select query"
	failedToExecuteDeleteQuery = "failed to execute a delete query"
	failedToExecuteUpdateQuery = "failed to execute a update query"
)

const (
	failedToScanRows        = "failed to scan rows"
	failedToIterateOverRows = "failed to iterate over rows"
	failedToScanRow         = "failed to scan rows"
)

const failedToGetDataFromDatabase = "failed to get data from database"

const (
	failedToCloseStatement = "failed to close statement"
	failedToCloseRows      = "failed to close rows"
)
