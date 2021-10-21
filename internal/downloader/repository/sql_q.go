package repository

var (
	pgGetAllFromTable = func(tableName string) string {
		return `select * from ` + tableName
	}
)
