package http

import (
	"github.com/labstack/echo/v4"
	"goland-hello/internal/tasks"
)

func NewDownloaderRoutes(dwGroup *echo.Group, h tasks.Handler) {
	dwGroup.POST("/create", h.Create())
	dwGroup.POST("/list", h.List())
	dwGroup.POST("/:emp_id", h.GetByEmployeeId())

	dwGroup.GET("/:tsk_id", h.GetByTaskId())
	dwGroup.PUT("/:tsk_id", h.Update())
	dwGroup.DELETE("/by_task/:tsk_id", h.DeleteByTaskId())
	dwGroup.DELETE("/by_employee/:emp_id", h.DeleteByEmployeeId())
}

