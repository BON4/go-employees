package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"goland-hello/internal/employees"
	"goland-hello/internal/models"
)

type EmployeeUC struct {
	repo employees.EmpRepository
	logger *logrus.Logger
}

func (e *EmployeeUC) GetByID(ctx context.Context, empID uint) (*models.Employee, error) {
	return e.repo.GetByID(ctx, empID)
}

func (e *EmployeeUC) Create(ctx context.Context, emp *models.Employee) (*models.Employee, error) {
	return e.repo.Create(ctx, emp)
}

func (e *EmployeeUC) Update(ctx context.Context, emp *models.Employee) (*models.Employee, error) {
	return e.repo.Update(ctx, emp)
}

func (e *EmployeeUC) Delete(ctx context.Context, empID uint) error {
	panic("implement me")
}

func (e *EmployeeUC) List(ctx context.Context, req *models.ListEmpRequest, dest []models.Employee) (int, error) {
	return e.repo.List(ctx, req, dest)
}

func NewEmployeeUC(repo employees.EmpRepository, logger *logrus.Logger) employees.EmpUC {
	return &EmployeeUC{
		repo:   repo,
		logger: logger,
	}
}