package server

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	e *echo.Echo
	cfg *config.Config
	db *pgxpool.Pool
	logger *logrus.Logger

	empFct models.EmployeeFactory
	tskFct models.TaskFactory
}

func NewServer(cfg *config.Config,
			   db *pgxpool.Pool,
			   empFct models.EmployeeFactory,
			   tskFct models.TaskFactory,
			   logger *logrus.Logger) *Server {
	return &Server{
		e:      echo.New(),
		cfg:    cfg,
		db:     db,
		empFct: empFct,
		tskFct: tskFct,
		logger: logger,
	}
}

func (s *Server) Run() error {
	if s.cfg.Server.SSL {
		return errors.New("SSL is not implemented")
	}

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.e.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(s.e); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.e.Server.Shutdown(ctx)
}
