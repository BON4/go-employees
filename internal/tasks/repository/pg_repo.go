package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	gerrors "github.com/pkg/errors"
	"goland-hello/internal/models"
	"goland-hello/internal/tasks"
	"goland-hello/pkg/dbErrors"
)

type tskPostgresRepo struct {
	taskTableName, empTableName string
	conn *pgxpool.Pool
}


func (t *tskPostgresRepo) Create(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	q := pgCreateTask(t.taskTableName)

	var createdTask models.Task
	err := t.conn.QueryRow(ctx, q, tsk.Open, tsk.Close, tsk.Closed,tsk.Meta, tsk.EmpId).
		Scan(
			&createdTask.TskId,
			&createdTask.Open,
			&createdTask.Close,
			&createdTask.Closed,
			&createdTask.Meta,
			&createdTask.EmpId,
		)

	var perr *pgconn.PgError
	if err != nil {
		//TODO Check errors codes end correct the dbErrors pkg
		if errors.As(err, &perr) {
			if perr.SQLState() == "23503" {
				return nil, dbErrors.NewViolates(perr, "tskPostgresRepo.Create: employee with this id does not exists, cant assignee task to non existing employee")
			}

			if perr.SQLState() == "23505" {
				return nil, dbErrors.NewAlreadyExists(err,"tskPostgresRepo.Create: task with this credentials already exists")
			}
		}
		return nil, gerrors.Wrap(err, "tskPostgresRepo.Create")
	}

	return &createdTask, err
}

func (t *tskPostgresRepo) Update(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	q := pgUpdateTask(t.taskTableName)

	var updatedTask models.Task

	err := t.conn.QueryRow(ctx, q, tsk.Open, tsk.Close, tsk.Closed,tsk.Meta, tsk.EmpId, tsk.TskId).
		Scan(
			&updatedTask.TskId,
			&updatedTask.Open,
			&updatedTask.Close,
			&updatedTask.Closed,
			&updatedTask.Meta,
			&updatedTask.EmpId,
		)

	var perr *pgconn.PgError
	if err != nil {
		//TODO Check errors codes end correct the dbErrors pkg
		if errors.As(err, &perr) {
			if perr.SQLState() == "23503" {
				return nil, dbErrors.NewViolates(perr, "tskPostgresRepo.Update: employee with this id does not exists, cant assignee task to non-existing employee")
			}

			if err == pgx.ErrNoRows {
				return nil, dbErrors.NewDoesNotExists(err,"tskPostgresRepo.Update: task with this credentials does not exists")
			}
		}
		return nil, gerrors.Wrap(err, "tskPostgresRepo.Update")
	}

	return &updatedTask, err
}

func (t *tskPostgresRepo) GetByTaskId(ctx context.Context, tskID uint) (*models.Task, error) {
	q := pgGetTaskByID(t.taskTableName)
	var foundTask models.Task

	err := t.conn.QueryRow(ctx, q, tskID).
		Scan(
			&foundTask.TskId,
			&foundTask.Open,
			&foundTask.Close,
			&foundTask.Closed,
			&foundTask.Meta,
			&foundTask.EmpId,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, dbErrors.NewDoesNotExists(err,"tskPostgresRepo.GetByID: task with this credentials does not exists")
		}
		return nil, gerrors.Wrap(err, "tskPostgresRepo.GetByID")
	}

	return &foundTask, nil
}

func (t *tskPostgresRepo) DeleteByTaskId(ctx context.Context, tskID uint) error {
	q := pgDeleteTask(t.taskTableName)
	ctag, err := t.conn.Exec(ctx, q, tskID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dbErrors.NewDoesNotExists(err,"tskPostgresRepo.Delete: task with this credentials does not exists")
		}
		return gerrors.Wrap(err, "tskPostgresRepo.Delete")
	}

	if ctag.RowsAffected() == 0 {
		return dbErrors.NewUnknown(errors.New("tskPostgresRepo.Delete"), "nothing has been deleted")
	}
	return nil
}

func (t *tskPostgresRepo) DeleteByEmployeeId(ctx context.Context, empId uint) error {
	q := pgDeleteTaskByEmpId(t.taskTableName)
	_, err := t.conn.Exec(ctx, q, empId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dbErrors.NewDoesNotExists(err,"tskPostgresRepo.Delete: tasks with this credentials does not exists")
		}
		return gerrors.Wrap(err, "tskPostgresRepo.Delete")
	}

	//if ctag.RowsAffected() == 0 {
	//	return dbErrors.NewUnknown(errors.New("tskPostgresRepo.Delete"), "nothing has been deleted")
	//}
	return nil
}

func (t *tskPostgresRepo) List(ctx context.Context, req *models.ListTskRequest, dest []models.Task) (int, error) {
	if len(dest) == 0 {
		return 0, nil
	}

	if req.PageSize == 0 {
		return 0, nil
	}

	q := pgListTask(t.taskTableName)
	rows, err := t.conn.Query(ctx, q, req.PageSize, req.PageNumber)
	if err != nil {
		return 0, gerrors.Wrap(err, "tskPostgresRepo.List")
	}
	defer rows.Close()

	var foundTask models.Task
	i := 0
	for rows.Next() {
		if ctx.Err() != nil {
			return i, gerrors.Wrap(ctx.Err(), "tskPostgresRepo.List.CtxErr")
		}

		if i >= len(dest) {
			return i, nil
		}

		err := rows.
			Scan(
			&foundTask.TskId,
			&foundTask.Open,
			&foundTask.Close,
			&foundTask.Closed,
			&foundTask.Meta,
			&foundTask.EmpId,
		)

		if err != nil {
			return 0, gerrors.Wrap(err, "tskPostgresRepo.List.Scan")
		}

		dest[i] = foundTask
		i++
	}

	if rows.Err() != nil {
		return 0, gerrors.Wrap(rows.Err(), "tskPostgresRepo.List.Rows")
	}

	return i, nil
}

func (t *tskPostgresRepo) GetByEmployeeId(ctx context.Context, empId uint, req *models.ListTskRequest, dest []models.Task) (int, error) {
	if len(dest) == 0 {
		return 0, nil
	}

	if req.PageSize == 0 {
		return 0, nil
	}

	q := pgGetTaskByEmpId(t.taskTableName, t.empTableName)
	rows, err := t.conn.Query(ctx, q, empId, req.PageSize, req.PageNumber)
	if err != nil {
		return 0, gerrors.Wrap(err, "tskPostgresRepo.GetByEmployeeId")
	}
	defer rows.Close()

	var foundTask models.Task
	i := 0
	for rows.Next() {
		if ctx.Err() != nil {
			return i, gerrors.Wrap(ctx.Err(), "tskPostgresRepo.GetByEmployeeId.CtxErr")
		}

		if i >= len(dest) {
			return i, nil
		}

		err := rows.
			Scan(
				&foundTask.TskId,
				&foundTask.Open,
				&foundTask.Close,
				&foundTask.Closed,
				&foundTask.Meta,
				&foundTask.EmpId,
			)

		if err != nil {
			return 0, gerrors.Wrap(err, "tskPostgresRepo.GetByEmployeeId.Scan")
		}

		dest[i] = foundTask
		i++
	}

	if rows.Err() != nil {
		return 0, gerrors.Wrap(rows.Err(), "tskPostgresRepo.GetByEmployeeId.Rows")
	}

	return i, nil
}

func NewTaskPostgresRepo(conn *pgxpool.Pool, taskTableName, empTableName string) tasks.TaskRepository {
	return &tskPostgresRepo{
		taskTableName: taskTableName,
		empTableName: empTableName,
		conn:      conn,
	}
}