package repository

var (
	pgCreateEmployee = func(tableName string) string {
		return `insert into `+tableName+` (fname, lname, sal) values ($1, $2, $3) returning *`
	}

	pgUpdateEmployee = func(tableName string) string {
		return `update `+tableName+` set fname = $1, lname = $2, sal = $3 where emp_id = $4 returning *`
	}

	pgDeleteEmployee = func(tableName string) string {
		return `delete from `+tableName+` where emp_id = $1`
	}

	pgListEmployee = func(tableName string) string {
		return `select * from `+tableName+` limit $1 offset $2 `
	}

	pgGetEmployeeByID = func(tableName string) string {
		return `select * from `+tableName+` where emp_id = $1`
	}
)