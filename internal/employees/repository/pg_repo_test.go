package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"goland-hello/config"
	"goland-hello/internal/models"
	"goland-hello/pkg/dbErrors"
	"testing"
)

var (
	skipDatabaseTest bool = false
	EmployeeFactory models.EmployeeFactory
	ConnDB *pgxpool.Pool
	ConfigDB config.Config
)

func flushEmployeeTable() error {
	_, err := ConnDB.Exec(context.Background(), "delete from employee")
	return err
}

func TestMain(m *testing.M) {
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

	if err = flushEmployeeTable(); err != nil {
		panic(err)
	}
}

func TestEmpPostgresRepo_Create(t *testing.T) {
	t.Parallel()

	repo := NewEmpPostgresRepo(ConnDB, "employee")
	t.Run("OK", func(t *testing.T) {
		empToCreate, err := EmployeeFactory.NewUser("test", "test", 120)
		require.Nil(t, err)

		createdEmp, err := repo.Create(context.Background(), &empToCreate)
		require.Nil(t, err)
		require.NotNil(t, createdEmp)
	})
}

func TestEmpPostgresRepo_GetByID(t *testing.T) {
	t.Parallel()
	repo := NewEmpPostgresRepo(ConnDB, "employee")
	t.Run("OK", func(t *testing.T) {
		empToCreate, err := EmployeeFactory.NewUser("test", "test", 120)
		require.Nil(t, err)

		createdEmp, err := repo.Create(context.Background(), &empToCreate)
		require.Nil(t, err)
		require.NotNil(t, createdEmp)

		foundEmp, err := repo.GetByID(context.Background(), createdEmp.EmpId)
		require.Nil(t, err)
		require.NotNil(t, foundEmp)

		require.Equal(t, createdEmp, foundEmp)
	})

	t.Run("FAIL", func(t *testing.T) {
		foundEmp, err := repo.GetByID(context.Background(), 0)
		require.Nil(t, foundEmp)
		require.NotNil(t, err)

		derr, ok := err.(dbErrors.DbErr)
		require.True(t, ok)
		require.NotNil(t, derr)
		require.Equal(t, dbErrors.ErrDoesNotExists, derr.Code())
	})
}

func TestEmpPostgresRepo_Delete(t *testing.T) {
	t.Parallel()
	repo := NewEmpPostgresRepo(ConnDB, "employee")
	t.Run("OK", func(t *testing.T) {
		empToCreate, err := EmployeeFactory.NewUser("test", "test", 120)
		require.Nil(t, err)

		createdEmp, err := repo.Create(context.Background(), &empToCreate)
		require.Nil(t, err)
		require.NotNil(t, createdEmp)

		err = repo.Delete(context.Background(), createdEmp.EmpId)
		require.Nil(t, err)
	})

	t.Run("FAIL", func(t *testing.T) {
		err := repo.Delete(context.Background(), 0)
		require.NotNil(t, err)

		derr, ok := err.(dbErrors.DbErr)
		require.True(t, ok)
		require.NotNil(t, derr)
		require.Equal(t, dbErrors.ErrorUnknown, derr.Code())
	})
}

func TestEmpPostgresRepo_Update(t *testing.T) {
	t.Parallel()
	repo := NewEmpPostgresRepo(ConnDB, "employee")
	t.Run("OK", func(t *testing.T) {
		empToCreate, err := EmployeeFactory.NewUser("test", "test", 120)
		require.Nil(t, err)

		createdEmp, err := repo.Create(context.Background(), &empToCreate)
		require.Nil(t, err)
		createdEmp.Sal = 1200
		updatedEmp, err := repo.Update(context.Background(), createdEmp)
		require.Nil(t, err)
		require.Equal(t, createdEmp.Sal, updatedEmp.Sal)
	})

	t.Run("FAIL", func(t *testing.T) {
		notExisting, err := EmployeeFactory.NewUser("fail_test", "Fail_test", 120)
		require.Nil(t, err)
		updatedEmp, err := repo.Update(context.Background(), &notExisting)
		require.Nil(t, updatedEmp)
		require.NotNil(t, err)

		derr, ok := err.(dbErrors.DbErr)
		require.True(t, ok)
		require.NotNil(t, derr)
	})
}

func TestEmpPostgresRepo_List(t *testing.T) {
	t.Parallel()
	repo := NewEmpPostgresRepo(ConnDB, "employee")
	empToCreate1, err := EmployeeFactory.NewUser("test", "test", 120)
	require.Nil(t, err)

	empToCreate2, err := EmployeeFactory.NewUser("test", "test", 120)
	require.Nil(t, err)

	_, err = repo.Create(context.Background(), &empToCreate1)
	require.Nil(t, err)

	_, err = repo.Create(context.Background(), &empToCreate2)
	require.Nil(t, err)


	dest := make([]models.Employee, 10)
	n, err := repo.List(context.Background(), &models.ListEmpRequest{
		PageSize:   10,
		PageNumber: 0,
	}, dest)
	require.Nil(t, err)
	require.Greater(t, n, 2)
	require.Equal(t, n, 3)
}