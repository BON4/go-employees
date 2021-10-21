package downloader

import (
	"context"
	"io"
)

type DWRepository interface {
	WriteTasks(ctx context.Context, writer io.Writer) (int, error)
	WriteEmployees(ctx context.Context, writer io.Writer) (int, error)

	GetHashTasks(ctx context.Context) (string, error)
	GetHashEmployees(ctx context.Context) (string, error)
}
