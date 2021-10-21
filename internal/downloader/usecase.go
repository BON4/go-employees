package downloader

import (
	"context"
)

type DwlUC interface {
	WriteTasks(ctx context.Context) (string, error)
	WriteEmployees(ctx context.Context) (string, error)
}
