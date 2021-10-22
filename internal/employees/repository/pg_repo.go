package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	gerrors "github.com/pkg/errors"
	"goland-hello/internal/employees"
	"goland-hello/internal/models"
	"goland-hello/pkg/dbErrors"
)

type empPostgresRepo struct {
	tableName string
	conn *pgxpool.Pool
}

func (e *empPostgresRepo) GetByID(ctx context.Context, empID uint) (*models.Employee, error) {
	q := pgGetEmployeeByID(e.tableName)
	var foundEmp models.Employee

	err := e.conn.QueryRow(ctx, q, empID).
		Scan(
			&foundEmp.EmpId,
			&foundEmp.Fname,
			&foundEmp.Lname,
			&foundEmp.Sal,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, dbErrors.NewDoesNotExists(err,"empPostgresRepo.GetByID: employee with this credentials does not exists")
		}
		return nil, gerrors.Wrap(err, "empPostgresRepo.GetByID")
	}

	return &foundEmp, nil
}

func (e *empPostgresRepo) Create(ctx context.Context, emp *models.Employee) (*models.Employee, error) {
	var perr *pgconn.PgError
	q := pgCreateEmployee(e.tableName)

	var createdEmp models.Employee

	err := e.conn.QueryRow(ctx, q, emp.Fname, emp.Lname, emp.Sal).
		Scan(
			&createdEmp.EmpId,
			&createdEmp.Fname,
			&createdEmp.Lname,
			&createdEmp.Sal,
			)

	if err != nil {
		//TODO Check errors codes end correct the dbErrors pkg
		if errors.As(err, &perr) {
			if perr.SQLState() == "23505" {
				return nil, dbErrors.NewAlreadyExists(err,"empPostgresRepo.Create: employee with this credentials already exists")
			}
		}
		return nil, gerrors.Wrap(err, "empPostgresRepo.Create")
	}

	return &createdEmp, nil
}

func (e *empPostgresRepo) Update(ctx context.Context, emp *models.Employee) (*models.Employee, error) {
	q := pgUpdateEmployee(e.tableName)

	var updatedEmp models.Employee
	err := e.conn.QueryRow(ctx, q, emp.Fname, emp.Lname, emp.Sal, emp.EmpId).
		Scan(
			&updatedEmp.EmpId,
			&updatedEmp.Fname,
			&updatedEmp.Lname,
			&updatedEmp.Sal,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, dbErrors.NewDoesNotExists(err,"empPostgresRepo.Update: employee with this credentials does not exists")
		}
		return nil, gerrors.Wrap(err, "empPostgresRepo.Update")
	}

	return &updatedEmp, nil
}

func (e *empPostgresRepo) Delete(ctx context.Context, empID uint) error {
	q := pgDeleteEmployee(e.tableName)
	ctag, err := e.conn.Exec(ctx, q, empID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dbErrors.NewDoesNotExists(err,"empPostgresRepo.Delete: employee with this credentials does not exists")
		}
		return gerrors.Wrap(err, "empPostgresRepo.Delete")
	}

	if ctag.RowsAffected() == 0 {
		return dbErrors.NewUnknown(errors.New("empPostgresRepo.Delete"), "nothing has been deleted")
	}
	return nil
}

func (e *empPostgresRepo) List(ctx context.Context, req *models.ListEmpRequest, dest []models.Employee) (int, error) {
	if len(dest) == 0 {
		return 0, nil
	}

	if req.PageSize == 0 {
		return 0, nil
	}

	q := pgListEmployee(e.tableName)
	rows, err := e.conn.Query(ctx, q, req.PageSize, req.PageNumber)
	if err != nil {
		return 0, gerrors.Wrap(err, "empPostgresRepo.List.Query")
	}
	defer rows.Close()

	var foundEmp models.Employee
	i := 0
	for rows.Next() {
		if ctx.Err() != nil {
			return 0, gerrors.Wrap(ctx.Err(), "empPostgresRepo.List.CtxErr")
		}

		if i >= len(dest) {
			return i, nil
		}

		err := rows.Scan(&foundEmp.EmpId, &foundEmp.Fname, &foundEmp.Lname, &foundEmp.Sal)
		if err != nil {
			return 0, gerrors.Wrap(err, "empPostgresRepo.List.Scan")
		}

		dest[i] = foundEmp
		i++
	}

	if rows.Err() != nil {
		return 0, gerrors.Wrap(rows.Err(), "empPostgresRepo.List.Rows")
	}

	return i, nil
}

func NewEmpPostgresRepo(conn *pgxpool.Pool, tableName string) employees.EmpRepository {
	return &empPostgresRepo{conn: conn, tableName: tableName}
}
