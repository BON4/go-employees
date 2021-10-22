package http

import (
	"github.com/labstack/echo/v4"
	"goland-hello/internal/downloader"
)

func NewDownloaderRoutes(dwGroup *echo.Group, h downloader.Handler) {
	dwGroup.GET("/employees", h.GetEmployees())
	dwGroup.GET("/tasks", h.GetTasks())
}