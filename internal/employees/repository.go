package employees

import (
	"context"
	"goland-hello/internal/models"
)

type EmpRepository interface {
	Create(ctx context.Context, emp *models.Employee) (*models.Employee, error)
	Update(ctx context.Context, emp *models.Employee) (*models.Employee, error)
	GetByID(ctx context.Context, empID uint) 		  (*models.Employee, error)
	Delete(ctx context.Context, empID uint)           error
	List(ctx context.Context, req *models.ListEmpRequest, dest []models.Employee) (int, error)
}
