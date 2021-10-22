package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/models"
	"goland-hello/internal/tasks"
	"goland-hello/pkg/httpErrors"
	"goland-hello/pkg/httpUtils"
	"net/http"
	"strconv"
)

type taskHandler struct {
	tskUC tasks.TaskUC
	tskFct models.TaskFactory
	logger *logrus.Logger
	cfg *config.Config
}

func (t *taskHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		tsk := &models.Task{}
		if err := httpUtils.ReadRequest(c, tsk); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		valTsk, err := t.tskFct.NewTask(tsk.EmpId, tsk.Open, tsk.Close, tsk.Closed, tsk.Meta)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdTsk, err := t.tskUC.Create(c.Request().Context(), &valTsk)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, createdTsk)
	}
}

func (t *taskHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		tsk := &models.Task{}
		if err := httpUtils.ReadRequest(c, tsk); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		valTsk, err := t.tskFct.Validate(tsk)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedTsk, err := t.tskUC.Update(c.Request().Context(), valTsk)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, updatedTsk)
	}
}

func (t *taskHandler) GetByTaskId() echo.HandlerFunc {
	return func(c echo.Context) error {
		tskId, err := strconv.ParseUint(c.Param("tsk_id"), 10, 32)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		foundTsk, err :=t.tskUC.GetByTaskId(c.Request().Context(), uint(tskId))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, foundTsk)
	}
}

func (t *taskHandler) DeleteByTaskId() echo.HandlerFunc {
	return func(c echo.Context) error {
		tskId, err := strconv.ParseUint(c.Param("tsk_id"), 10, 32)

		err = t.tskUC.DeleteByTaskId(c.Request().Context(), uint(tskId))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.NoContent(http.StatusOK)
	}
}

func (t *taskHandler) DeleteByEmployeeId() echo.HandlerFunc {
	return func(c echo.Context) error {
		empId, err := strconv.ParseUint(c.Param("emp_id"), 10, 32)

		err = t.tskUC.DeleteByEmployeeId(c.Request().Context(), uint(empId))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.NoContent(http.StatusOK)
	}
}

func (t *taskHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		tskReq := &models.ListTskRequest{}
		if err := httpUtils.ReadRequest(c, tskReq); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//TODO later, could be optimized with sync.Pool
		dest := make([]models.Task, tskReq.PageSize)

		n, err := t.tskUC.List(c.Request().Context(), tskReq, dest)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, dest[:n])
	}
}

func (t *taskHandler) GetByEmployeeId() echo.HandlerFunc {
	return func(c echo.Context) error {
		empId, err := strconv.ParseUint(c.Param("emp_id"), 10, 32)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		tskReq := &models.ListTskRequest{}
		if err := httpUtils.ReadRequest(c, tskReq); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//TODO later, could be optimized with sync.Pool
		dest := make([]models.Task, tskReq.PageSize)

		n, err := t.tskUC.GetByEmployeeId(c.Request().Context(), uint(empId), tskReq, dest)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, dest[:n])
	}
}

func NewTaskHandler(tskUC tasks.TaskUC,
					tskFct models.TaskFactory,
					logger *logrus.Logger,
					cfg *config.Config) tasks.Handler {
	return &taskHandler{
		tskFct: tskFct,
		tskUC:  tskUC,
		logger: logger,
		cfg:    cfg,
	}
}