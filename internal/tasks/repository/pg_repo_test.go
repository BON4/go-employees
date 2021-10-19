package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"goland-hello/config"
	"goland-hello/internal/employees/repository"
	"goland-hello/internal/employees/usecase"
	"goland-hello/internal/models"
	"goland-hello/pkg/dbErrors"
	"testing"
	"time"
)

var (
	skipDatabaseTest bool = false
	TaskFactory models.TaskFactory
	EmployeeFactory models.EmployeeFactory
	ConnDB *pgxpool.Pool
	ConfigDB config.Config
)

func flushEmployeeTable() error {
	_, err := ConnDB.Exec(context.Background(), "delete from employee")
	return err
}

func flushTaskTable() error {
	_, err := ConnDB.Exec(context.Background(), "delete from task")
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

func TestTskPostgresRepo_Create(t *testing.T) {
	t.Parallel()
	tskRepo := NewTaskPostgresRepo(ConnDB, "task")
	empRepo := repository.NewEmpPostgresRepo(ConnDB, "employee")
	empUC := usecase.NewEmployeeUC(empRepo, nil)

	t.Run("OK", func(t *testing.T) {
		fEmp, err := EmployeeFactory.NewUser("test", "test", 120)
		require.NoError(t, err)

		createdEmp, err := empUC.Create(context.Background(), &fEmp)
		require.NoError(t, err)

		task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, "")
		require.NoError(t, err)

		createdTask, err := tskRepo.Create(context.Background(), &task)
		require.NoError(t, err)
		require.NotNil(t, createdTask)
	})

	t.Run("FAIL No employee", func(t *testing.T) {
		task, err := TaskFactory.NewTask(0, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, "")
		require.NoError(t, err)

		createdTask, err := tskRepo.Create(context.Background(), &task)
		require.Error(t, err)
		require.Nil(t, createdTask)

		derr, ok := err.(dbErrors.DbErr)
		require.True(t, ok)
		require.NotNil(t, derr)
		require.Equal(t, dbErrors.ErrorViolates, derr.Code())
	})
}

func TestTskPostgresRepo_Update(t *testing.T) {
	t.Parallel()
	tskRepo := NewTaskPostgresRepo(ConnDB, "task")
	empRepo := repository.NewEmpPostgresRepo(ConnDB, "employee")
	empUC := usecase.NewEmployeeUC(empRepo, nil)

	t.Run("OK", func(t *testing.T) {
		fEmp, err := EmployeeFactory.NewUser("test", "test", 120)
		require.NoError(t, err)

		createdEmp, err := empUC.Create(context.Background(), &fEmp)
		require.NoError(t, err)

		task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, "")
		require.NoError(t, err)

		createdTask, err := tskRepo.Create(context.Background(), &task)
		require.NoError(t, err)
		require.NotNil(t, createdTask)

		newTask, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*15).Unix(), false, "")
		require.NoError(t, err)
		newTask.TskId = createdTask.TskId

		updatedTask, err := tskRepo.Update(context.Background(), &newTask)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)

		require.Equal(t, newTask, *updatedTask)
	})

	t.Run("FAIL No employee", func(t *testing.T) {
		fEmp, err := EmployeeFactory.NewUser("test", "test", 120)
		require.NoError(t, err)

		createdEmp, err := empUC.Create(context.Background(), &fEmp)
		require.NoError(t, err)

		task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, "")
		require.NoError(t, err)

		createdTask, err := tskRepo.Create(context.Background(), &task)
		require.NoError(t, err)
		require.NotNil(t, createdTask)
		createdTask.EmpId = 0

		updatedTask, err := tskRepo.Update(context.Background(), createdTask)
		require.Error(t, err)
		require.Nil(t, updatedTask)

		derr, ok := err.(dbErrors.DbErr)
		require.True(t, ok)
		require.NotNil(t, derr)
		require.Equal(t, dbErrors.ErrorViolates, derr.Code())
	})
}