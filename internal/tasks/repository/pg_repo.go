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
	tableName string
	conn *pgxpool.Pool
}

func (t tskPostgresRepo) Create(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	q := pgCreateTask(t.tableName)

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

func (t tskPostgresRepo) Update(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	q := pgUpdateTask(t.tableName)

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

func (t tskPostgresRepo) GetByID(ctx context.Context, tskID uint) (*models.Task, error) {
	panic("implement me")
}

func (t tskPostgresRepo) Delete(ctx context.Context, tskID uint) error {
	panic("implement me")
}

func (t tskPostgresRepo) List(ctx context.Context, req *models.ListTskRequest, dest []models.Task) (int, error) {
	panic("implement me")
}

func NewTaskPostgresRepo(conn *pgxpool.Pool, tableName string) tasks.TaskRepository {
	return &tskPostgresRepo{
		tableName: tableName,
		conn:      conn,
	}
}