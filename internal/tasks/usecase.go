package tasks

import (
	"context"
	"goland-hello/internal/models"
)

type TaskUC interface {
	Create(ctx context.Context, tsk *models.Task) (*models.Task, error)
	Update(ctx context.Context, tsk *models.Task) (*models.Task, error)
	GetByID(ctx context.Context, tskID uint) 	  (*models.Task, error)
	DeleteByTaskId(ctx context.Context, tskID uint) error
	DeleteByEmployeeId(ctx context.Context, empId uint) error
	List(ctx context.Context, req *models.ListTskRequest, dest []models.Task) (int, error)
}
