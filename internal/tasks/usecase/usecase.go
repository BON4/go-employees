package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"goland-hello/internal/models"
	"goland-hello/internal/tasks"
)

type TskUC struct {
	tskRepo tasks.TaskRepository
	logger *logrus.Logger
}

func (t *TskUC) Create(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	return t.tskRepo.Create(ctx, tsk)
}

func (t *TskUC) Update(ctx context.Context, tsk *models.Task) (*models.Task, error) {
	return t.tskRepo.Update(ctx, tsk)
}

func (t *TskUC) GetByTaskId(ctx context.Context, tskID uint) (*models.Task, error) {
	return t.tskRepo.GetByTaskId(ctx, tskID)
}

func (t *TskUC) DeleteByTaskId(ctx context.Context, tskID uint) error {
	return t.tskRepo.DeleteByTaskId(ctx, tskID)
}

func (t *TskUC) DeleteByEmployeeId(ctx context.Context, empId uint) error {
	return t.tskRepo.DeleteByEmployeeId(ctx, empId)
}

func (t *TskUC) List(ctx context.Context, req *models.ListTskRequest, dest []models.Task) (int, error) {
	return t.tskRepo.List(ctx, req, dest)
}

func (t *TskUC) GetByEmployeeId(ctx context.Context, empId uint, req *models.ListTskRequest, dest []models.Task) (int, error) {
	panic("implement me")
}

func NewTaskUseCase(tskRepo tasks.TaskRepository, logger *logrus.Logger) tasks.TaskUC {
	return &TskUC{
		tskRepo: tskRepo,
		logger:  logger,
	}
}
