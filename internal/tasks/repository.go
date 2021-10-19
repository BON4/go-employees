package tasks

import (
	"context"
	"goland-hello/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, tsk *models.Task) (*models.Task, error)
	Update(ctx context.Context, tsk *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, tskID uint) 	  (*models.Task, error)
	Delete(ctx context.Context, tskID uint)       error
	List(ctx context.Context, req *models.ListTskRequest, dest []models.Task) (int, error)
}
