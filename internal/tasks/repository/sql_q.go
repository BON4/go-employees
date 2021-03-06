package repository

var (
	pgCreateTask = func(tableName string) string {
		return `insert into `+tableName+` (open_d, close_d, closed ,meta, emp_id) values ($1, $2, $3, $4, $5) returning *`
	}

	pgUpdateTask = func(tableName string) string {
		return `update `+tableName+` set open_d = $1, close_d = $2,closed = $3, meta = $4, emp_id = $5 where tsk_id = $6 returning *`
	}

	pgDeleteTask = func(tableName string) string {
		return `delete from `+tableName+` where tsk_id = $1`
	}

	pgDeleteTaskByEmpId = func(tableName string) string {
		return `delete from `+tableName+` where emp_id = $1`
	}

	pgListTask = func(tableName string) string {
		return `select * from `+tableName+` limit $1 offset $2 `
	}

	pgGetTaskByID = func(tableName string) string {
		return `select * from `+tableName+` where tsk_id = $1`
	}

	pgGetTaskByEmpId = func(taskTableName string, empTableName string) string {
		return `select t.* from ` + empTableName +` inner join ` + taskTableName + ` t on employee.emp_id = t.emp_id where t.emp_id = $1 limit $2 offset $3`
	}
)
