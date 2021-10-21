package repository

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"goland-hello/config"
	empRepoMd "goland-hello/internal/employees/repository"
	empUcMd "goland-hello/internal/employees/usecase"
	"goland-hello/internal/models"
	"goland-hello/internal/pkg/utils"
	tskRepoMd "goland-hello/internal/tasks/repository"
	tskUcMd "goland-hello/internal/tasks/usecase"
	"io"
	"sync"
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

// ! WARNING !
// These tests should not been executed parallel or concurrently
func TestDownloaderRepository_Write(t *testing.T) {
	require.NoError(t, flushTaskTable())
	require.NoError(t, flushEmployeeTable())

	tskRepo := tskRepoMd.NewTaskPostgresRepo(ConnDB, taskTableName, empTableName)
	tskUc := tskUcMd.NewTaskUseCase(tskRepo, nil)

	empRepo := empRepoMd.NewEmpPostgresRepo(ConnDB, empTableName)
	empUC := empUcMd.NewEmployeeUC(empRepo, nil)

	dwRepo := NewDownloaderRepository(ConnDB, taskTableName, empTableName)

	t.Run("WriteTask", func(t *testing.T) {
		emp, err := EmployeeFactory.NewUser("test", "test", 1200)
		require.NoError(t, err)

		createdEmp, err := empUC.Create(context.Background(), &emp)
		require.NoError(t, err)

		rActual, wActual := io.Pipe()
		rTest, wTest := io.Pipe()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup) {
			defer wg.Done()
			actualCsv, err := io.ReadAll(rActual)
			require.NoError(t, err)
			testCsv, err := io.ReadAll(rTest)
			require.NoError(t, err)
			require.Equal(t, 0, bytes.Compare(testCsv, actualCsv))
		}(t, wg)

		testBuf := bufio.NewWriter(wActual)
		for i := 0; i < 10; i++ {
			task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, fmt.Sprintf("Tsk #%d", i))
			require.NoError(t, err)

			createdTsk, err := tskUc.Create(context.Background(), &task)
			require.NoError(t, err)

			_, err = utils.TaskToByte(createdTsk, testBuf)
			require.NoError(t, err)
		}
		testBuf.Flush()
		wActual.Close()

		_, err = dwRepo.WriteTasks(context.Background(), wTest)
		require.NoError(t, err)
		wTest.Close()

		wg.Wait()
	})

	require.NoError(t, flushTaskTable())
	require.NoError(t, flushEmployeeTable())

	t.Run("WriteEmployees", func(t *testing.T) {
		rActual, wActual := io.Pipe()
		rTest, wTest := io.Pipe()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(t *testing.T, wg *sync.WaitGroup) {
			defer wg.Done()
			actualCsv, err := io.ReadAll(rActual)
			require.NoError(t, err)
			testCsv, err := io.ReadAll(rTest)
			require.NoError(t, err)
			require.Equal(t, 0, bytes.Compare(testCsv, actualCsv))
		}(t, wg)

		testBuf := bufio.NewWriter(wActual)
		for i := 0; i < 10; i++ {
			emp, err := EmployeeFactory.NewUser("test", "test", 1200)
			require.NoError(t, err)

			createdEmp, err := empUC.Create(context.Background(), &emp)
			require.NoError(t, err)

			_, err = utils.EmpToByte(createdEmp, testBuf)
			require.NoError(t, err)
		}
		testBuf.Flush()
		wActual.Close()

		_, err := dwRepo.WriteEmployees(context.Background(), wTest)
		require.NoError(t, err)
		wTest.Close()

		wg.Wait()
	})

	require.NoError(t, flushTaskTable())
	require.NoError(t, flushEmployeeTable())
}

func TestDownloaderRepository_GetHash(t *testing.T) {
	tskRepo := tskRepoMd.NewTaskPostgresRepo(ConnDB, taskTableName, empTableName)
	tskUc := tskUcMd.NewTaskUseCase(tskRepo, nil)

	empRepo := empRepoMd.NewEmpPostgresRepo(ConnDB, empTableName)
	empUC := empUcMd.NewEmployeeUC(empRepo, nil)

	dwRepo := NewDownloaderRepository(ConnDB, taskTableName, empTableName)

	t.Run("Emplyee table md5 hash", func(t *testing.T) {
		emp, err := EmployeeFactory.NewUser("test", "test", 1200)
		require.NoError(t, err)

		_, err = empUC.Create(context.Background(), &emp)
		require.NoError(t, err)

		hash, err := dwRepo.GetHashEmployees(context.Background())
		require.NoError(t, err)
		require.Equal(t, 32, len(hash))
		fmt.Printf("Emplyee hash: %q\n", hash)
	})

	t.Run("Employee table md5 hash", func(t *testing.T) {
		emp, err := EmployeeFactory.NewUser("test", "test", 1200)
		require.NoError(t, err)

		createdEmp, err := empUC.Create(context.Background(), &emp)
		require.NoError(t, err)

		task, err := TaskFactory.NewTask(createdEmp.EmpId, time.Now().Unix(), time.Now().Add(time.Hour*10).Unix(), false, "TASK 1")
		require.NoError(t, err)

		_, err = tskUc.Create(context.Background(), &task)
		require.NoError(t, err)

		hash, err := dwRepo.GetHashTasks(context.Background())
		require.NoError(t, err)
		require.Equal(t, 32, len(hash))
		fmt.Printf("Task hash: %q\n", hash)
	})
}