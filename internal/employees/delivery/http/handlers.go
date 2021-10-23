package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/employees"
	"goland-hello/internal/models"
	"goland-hello/pkg/httpErrors"
	"goland-hello/pkg/httpUtils"
	"net/http"
	"strconv"
)

type employeeHandler struct {
	empUC employees.EmpUC
	empFct models.EmployeeFactory
	logger *logrus.Logger
	cfg *config.Config
}

func (e *employeeHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		emp := &models.Employee{}
		if err := httpUtils.ReadRequest(c, emp); err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		valEmp, err := e.empFct.NewUser(emp.Fname, emp.Lname, emp.Sal)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdEmp, err := e.empUC.Create(c.Request().Context(), &valEmp)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdEmp)
	}
}

func (e *employeeHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		emp := &models.Employee{}
		if err := httpUtils.ReadRequest(c, emp); err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		valEmp, err := e.empFct.Validate(emp)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedEmp, err := e.empUC.Update(c.Request().Context(), valEmp)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, updatedEmp)
	}
}

func (e *employeeHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		empId, err := strconv.ParseUint(c.Param("emp_id"), 10, 32)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		err = e.empUC.Delete(c.Request().Context(), uint(empId))
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.NoContent(http.StatusOK)
	}
}

func (e *employeeHandler) GetById() echo.HandlerFunc {
	return func(c echo.Context) error {
		empId, err := strconv.ParseUint(c.Param("emp_id"), 10, 32)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		foundEmp, err := e.empUC.GetByID(c.Request().Context(), uint(empId))
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, foundEmp)
	}
}

func (e *employeeHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		empReq := &models.ListEmpRequest{}
		if err := httpUtils.ReadRequest(c, empReq); err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//TODO later, could be optimized with sync.Pool
		dest := make([]models.Employee, empReq.PageSize)

		n, err := e.empUC.List(c.Request().Context(), empReq, dest)
		if err != nil {
			log.Warn(err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, dest[:n])
	}
}

func NewEmployeeHandler(empUC employees.EmpUC,
						empFct models.EmployeeFactory,
						logger *logrus.Logger,
						cfg *config.Config) employees.Handler {
	return &employeeHandler{
		empUC:  empUC,
		empFct: empFct,
		logger: logger,
		cfg:    cfg,
	}
}