package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"goland-hello/config"
	empRepoMd "goland-hello/internal/employees/repository"
	empUcMd "goland-hello/internal/employees/usecase"
	"goland-hello/internal/models"
	tskRepoMd "goland-hello/internal/tasks/repository"
	tskUcMd "goland-hello/internal/tasks/usecase"
	"os"
	"testing"
	"time"
)

const (
	taskTableName = "task"
	empTableName = "employee"
)


var (
	skipDatabaseTest bool = false
	TaskFactory models.TaskFactory
	EmployeeFactory models.EmployeeFactory
	ConnDB *pgxpool.Pool
	ConfigDB config.Config
)

func flushEmployeeTable() error {
	_, err := ConnDB.Exec(context.Background(), "delete from "+empTableName)
	return err
}

func flushTaskTable() error {
	_, err := ConnDB.Exec(context.Background(), "delete from "+taskTableName)
	return err
}

func TestMain(m *testing.M) {
	TaskFactory = models.NewTaskFactory(models.TaskFactoryConfig{
		MinTaskLifespan: time.Hour*1,
	})

	EmployeeFactory = models.NewEmployeeFactory(models.EmployeeFactoryConfig{
		MinFirstNameLength: 3,
		MinLastNameLength:  3,
		MinSalary:          0,
	})

	ConfigDB, err := config.ParseConfig("/home/home/go/github.com/BON4/go-employees/config/test_conf.yaml")
	if err != nil {
		panic(err)
	}

	ConnDB, err = config.OpenPostgresPoolConfig(context.Background(), &ConfigDB)
	if err != nil {
		panic(err)
	}
	defer ConnDB.Close()

	err = ConnDB.Ping(context.Background())
	if err != nil {
		println("Cant connect")
		skipDatabaseTest = true
	} else {
		println("Database connected")
	}
	m.Run()

	errTsk, errEmp := flushTaskTable(), flushEmployeeTable()
	if errTsk != nil {
		panic(errTsk)
	} else if errEmp != nil {
		panic(errEmp)
	}
}

func TestTaskToByte(t *testing.T) {
	tskRepo := tskRepoMd.NewTaskPostgresRepo(ConnDB, taskTableName, empTableName)
	tskUc := tskUcMd.NewTaskUseCase(tskRepo, nil)

	empRepo := empRepoMd.NewEmpPostgresRepo(ConnDB, empTableName)
	empUC := empUcMd.NewEmployeeUC(empRepo, nil)

	dwRepo := NewDownloaderRepository(ConnDB)

	emp, err := EmployeeFactory.NewUser("test", "test", 1200)
	require.NoError(t, err)

	createdEmp, err := empUC.Create(context.Background(), &emp)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, fmt.Sprintf("Tsk #%d", i))
		require.NoError(t, err)
		_, err = tskUc.Create(context.Background(), &task)
		require.NoError(t, err)
	}

	_, err = dwRepo.WriteTasks(context.Background(), taskTableName, os.Stdout)
	require.NoError(t, err)
}