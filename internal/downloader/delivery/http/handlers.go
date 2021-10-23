package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/downloader"
	"goland-hello/pkg/httpErrors"
	"net/http"
)

type downloaderHandler struct {
	dwlUC downloader.DwlUC
	logger *logrus.Logger
	cfg *config.Config
}

func (d *downloaderHandler) GetTasks() echo.HandlerFunc {
	return func(c echo.Context) error {
		hash, err := d.dwlUC.WriteTasks(c.Request().Context())
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, hash)
	}
}

func (d *downloaderHandler) GetEmployees() echo.HandlerFunc {
	return func(c echo.Context) error {
		hash, err := d.dwlUC.WriteEmployees(c.Request().Context())
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, hash)
	}
}

func NewDownloaderHandler(dwlUC downloader.DwlUC,
						  logger *logrus.Logger,
						  cfg *config.Config,
						  ) downloader.Handler {
	return &downloaderHandler{
		dwlUC:      dwlUC,
		logger:     logger,
		cfg:        cfg,
	}
}
