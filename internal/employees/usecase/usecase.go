package usecase
//Use-Cases needed for combining business logic.
//For example, in use-case you can combine logic of saving data to multiple databases under one method

import (
	"context"
	"github.com/sirupsen/logrus"
	"goland-hello/internal/employees"
	"goland-hello/internal/models"
	"goland-hello/internal/tasks"
)


type EmployeeUC struct {
	repo employees.EmpRepository
	tskUC tasks.TaskUC
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
	//Example of using a use-case, delete from one table and another
	//Of course you can do it with cascade delete, but I implemented it only for educational proposes (sometimes you want this logic in service, not in db)
	err := e.tskUC.DeleteByEmployeeId(ctx, empID)
	if err != nil {
		return err
	}
	return e.repo.Delete(ctx, empID)
}

func (e *EmployeeUC) List(ctx context.Context, req *models.ListEmpRequest, dest []models.Employee) (int, error) {
	//Here you can easily add some caching, maybe TODO
	return e.repo.List(ctx, req, dest)
}

func NewEmployeeUC(repo employees.EmpRepository, tskUc tasks.TaskUC,logger *logrus.Logger) employees.EmpUC {
	return &EmployeeUC{
		repo:   repo,
		tskUC: tskUc,
		logger: logger,
	}
}