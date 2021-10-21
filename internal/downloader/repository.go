package downloader

import (
	"context"
	"io"
)

type DWRepository interface {
	WriteTasks(ctx context.Context, tableName string, writer io.Writer) (int, error)
	WriteEmployees(ctx context.Context, tableName string, writer io.Writer) (int, error)
}
