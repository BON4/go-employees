package repository

var (
	pgGetAllFromTable = func(tableName string) string {
		return `select * from ` + tableName
	}

	pgGetMD5HashOfEmployeeTable = func(tableName string) string {
		return `SELECT md5(CAST((array_agg(f.* order by emp_id))AS text)) FROM `+ tableName +` f`
	}

	pgGetMD5HashOfTaskTable = func(tableName string) string {
		return `SELECT md5(CAST((array_agg(f.* order by tsk_id))AS text)) FROM `+ tableName +` f`
	}
)
