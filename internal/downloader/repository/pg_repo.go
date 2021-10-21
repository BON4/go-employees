package repository

import (
	"bufio"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	gerrors "github.com/pkg/errors"
	"goland-hello/internal/downloader"
	"goland-hello/internal/models"
	"goland-hello/internal/pkg/utils"
	"io"
)

type DownloaderRepository struct {
	conn *pgxpool.Pool
}

//TODO SHITS NEEDS GENERICS FOR CLEANER IMPLEMENTATION

// WriteEmployees - Writes whole employee table to writer in csv RFC 4180 format, returns number of bytes that have been written and error
func (d *DownloaderRepository) WriteEmployees(ctx context.Context, tableName string, writer io.Writer) (int, error) {
	q := pgGetAllFromTable(tableName)
	rows, err := d.conn.Query(ctx, q)
	if err != nil {
		return 0, gerrors.Wrap(err, "DownloaderRepository.WriteEmployees")
	}
	defer rows.Close()

	bw := bufio.NewWriter(writer)

	n := 0
	var foundEmp models.Employee
	for rows.Next() {
		err := rows.Scan(&foundEmp.EmpId, &foundEmp.Fname, &foundEmp.Lname, &foundEmp.Sal)

		if err != nil {
			return 0, gerrors.Wrap(err, "DownloaderRepository.WriteEmployees")
		}

		if m, err := utils.EmpToByte(&foundEmp, bw); err != nil {
			return 0, gerrors.Wrap(err, "DownloaderRepository.WriteEmployees.TaskToByte")
		} else {
			n += m
		}
	}

	//TODO (later) buffer may overflow if table will be too large, Flush it maybe avery N bytes
	err = bw.Flush()
	if err != nil {
		return 0, gerrors.Wrap(err, "DownloaderRepository.WriteEmployees.Flush")
	}

	if rows.Err() != nil {
		return 0, gerrors.Wrap(rows.Err(), "DownloaderRepository.WriteEmployees.Rows")
	}
	return n, nil
}

// WriteTasks - Writes whole task table to writer in csv RFC 4180 format, returns number of bytes that have been written and error
func (d *DownloaderRepository) WriteTasks(ctx context.Context, tableName string, writer io.Writer) (int, error) {
	q := pgGetAllFromTable(tableName)
	rows, err := d.conn.Query(ctx, q)
	if err != nil {
		return 0, gerrors.Wrap(err, "DownloaderRepository.WriteTasks")
	}
	defer rows.Close()

	bw := bufio.NewWriter(writer)

	n := 0
	var foundTask models.Task
	for rows.Next() {
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
			return 0, gerrors.Wrap(err, "DownloaderRepository.WriteTasks")
		}

		if m, err := utils.TaskToByte(&foundTask, bw); err != nil {
			return 0, gerrors.Wrap(err, "DownloaderRepository.WriteTasks.TaskToByte")
		} else {
			n += m
		}
	}

	//TODO (later) buffer may overflow if table will be too large, Flush it maybe avery N bytes
	err = bw.Flush()
	if err != nil {
		return 0, gerrors.Wrap(err, "DownloaderRepository.WriteTasks.Flush")
	}

	if rows.Err() != nil {
		return 0, gerrors.Wrap(rows.Err(), "DownloaderRepository.WriteTasks.Rows")
	}
	return n, nil
}

func NewDownloaderRepository(conn *pgxpool.Pool) downloader.DWRepository {
	return &DownloaderRepository{conn: conn}
}