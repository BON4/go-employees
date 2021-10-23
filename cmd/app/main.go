package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/models"
	"goland-hello/internal/server"
	"time"
)

func main(){
	TaskFactory := models.NewTaskFactory(models.TaskFactoryConfig{
		MinTaskLifespan: time.Hour*1,
	})

	EmployeeFactory := models.NewEmployeeFactory(models.EmployeeFactoryConfig{
		MinFirstNameLength: 3,
		MinLastNameLength:  3,
		MinSalary:          0,
	})

	Config, err := config.ParseConfig("/home/home/go/github.com/BON4/go-employees/config/test_conf.yaml")
	ConnDB, err := config.OpenPostgresPoolConfig(context.Background(), &Config)
	if err != nil {
		panic(err)
	}
	defer ConnDB.Close()

	logger := logrus.New()

	s := server.NewServer(&Config, ConnDB, EmployeeFactory, TaskFactory, logger)
	if err = s.Run(); err != nil {
		logger.Fatal(err)
	}
}