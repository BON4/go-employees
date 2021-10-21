package usecase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"goland-hello/config"
	dwRepoMd "goland-hello/internal/downloader/repository"
	empRepoMd "goland-hello/internal/employees/repository"
	empUcMd "goland-hello/internal/employees/usecase"
	"goland-hello/internal/models"
	//tskRepoMd "goland-hello/internal/tasks/repository"
	//tskUcMd "goland-hello/internal/tasks/usecase"
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
	var err error
	ConfigDB, err = config.ParseConfig("/home/home/go/github.com/BON4/go-employees/config/test_conf.yaml")
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

func TestDownloaderUC_WriteEmployees(t *testing.T) {
	//tskRepo := tskRepoMd.NewTaskPostgresRepo(ConnDB, taskTableName, empTableName)
	//tskUc := tskUcMd.NewTaskUseCase(tskRepo, nil)

	empRepo := empRepoMd.NewEmpPostgresRepo(ConnDB, empTableName)
	empUC := empUcMd.NewEmployeeUC(empRepo, nil)

	dwRepo := dwRepoMd.NewDownloaderRepository(ConnDB, taskTableName, empTableName)
	dwUC := NewDownloaderUC(dwRepo, &ConfigDB, nil)

	emp, err := EmployeeFactory.NewUser("test", "test", 1200)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		_, err := empUC.Create(context.Background(), &emp)
		require.NoError(t, err)
	}

	fielname, err := dwUC.WriteEmployees(context.Background())
	require.NoError(t, err)
	fmt.Println(fielname)
}