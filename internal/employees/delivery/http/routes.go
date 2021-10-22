package http

import (
	"github.com/labstack/echo/v4"
	"goland-hello/internal/employees"
)

func NewDownloaderRoutes(dwGroup *echo.Group, h employees.Handler) {
	dwGroup.POST("/create", h.Create())
	dwGroup.POST("/list", h.List())

	dwGroup.GET("/:emp_id", h.GetById())
	dwGroup.PUT("/:emp_id", h.Update())
	dwGroup.DELETE("/:emp_id", h.Delete())
}
