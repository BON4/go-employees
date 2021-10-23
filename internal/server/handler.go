package server

import (
	"github.com/labstack/echo/v4"
	dwlHttp "goland-hello/internal/downloader/delivery/http"
	dwlRepoMd "goland-hello/internal/downloader/repository"
	dwlUcMd "goland-hello/internal/downloader/usecase"
	empHttp "goland-hello/internal/employees/delivery/http"
	empRepoMd "goland-hello/internal/employees/repository"
	empUcMd "goland-hello/internal/employees/usecase"
	tskHttp "goland-hello/internal/tasks/delivery/http"
	tskRepoMd "goland-hello/internal/tasks/repository"
	tskUcMd "goland-hello/internal/tasks/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	tskRepo := tskRepoMd.NewTaskPostgresRepo(s.db, "task", "employee")
	tskUc := tskUcMd.NewTaskUseCase(tskRepo, s.logger)
	tskRoute := tskHttp.NewTaskHandler(tskUc, s.tskFct, s.logger, s.cfg)

	empRepo := empRepoMd.NewEmpPostgresRepo(s.db, "employee")
	empUc := empUcMd.NewEmployeeUC(empRepo, tskUc,s.logger)
	empRoute := empHttp.NewEmployeeHandler(empUc, s.empFct, s.logger, s.cfg)

	dwlRepo := dwlRepoMd.NewDownloaderRepository(s.db, "task", "employee")
	dwlUc := dwlUcMd.NewDownloaderUC(dwlRepo, s.cfg, s.logger)
	dwlRoute := dwlHttp.NewDownloaderHandler(dwlUc, s.logger, s.cfg)

	v1 := e.Group("/v1")

	empGroup := v1.Group("/emp")
	tskGroup := v1.Group("/tsk")
	dwlGroup := v1.Group("/dwl")

	empHttp.NewEmployeeRoutes(empGroup, empRoute)
	tskHttp.NewTaskRoutes(tskGroup, tskRoute)
	dwlHttp.NewDownloaderRoutes(dwlGroup, dwlRoute)
	return nil
}
