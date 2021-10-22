package downloader

import "github.com/labstack/echo/v4"

type Handler interface {
	GetTasks()     echo.HandlerFunc
	GetEmployees() echo.HandlerFunc
}
