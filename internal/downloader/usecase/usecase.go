package usecase

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"goland-hello/config"
	"goland-hello/internal/downloader"
	"os"
)

type DownloaderUC struct {
	repo downloader.DWRepository
	cfg  config.Downloader
	logger *logrus.Logger
}

// WriteTasks - write task table to csv file, and return its name
// In case if table did not change, returns name of csv file that up to date
func (d *DownloaderUC) WriteTasks(ctx context.Context) (string, error) {
	hash, err := d.repo.GetHashTasks(ctx)
	if err != nil {
		return "", err
	}

	hashFile := d.cfg.FileFolder+"tsk_"+hash+".csv"
	if _, err := os.Stat(hashFile); os.IsNotExist(err) {
		// csv file not exists, this means that csv is not upto date

		//TODO dont know how this would behave with parallel
		f, err := os.Create(hashFile)
		if err != nil {
			return "", err
		}

		defer f.Close()

		_, err = d.repo.WriteTasks(ctx, f)
		if err != nil {
			return "", err
		}

		if err = f.Sync(); err != nil {
			return "", err
		}
	}

	return hashFile, nil
}

// WriteEmployees - write employee table to csv file, and return its name
// In case if table did not change, returns name of csv file that up to date
func (d *DownloaderUC) WriteEmployees(ctx context.Context) (string, error) {
	hash, err := d.repo.GetHashEmployees(ctx)
	if err != nil {
		return "", err
	}

	hashFile := d.cfg.FileFolder+"emp_"+hash+".csv"
	if _, err := os.Stat(hashFile); os.IsNotExist(err) {
		// csv file not exists, this means that csv is not upto date
		//TODO dont know how this would behave with parallel
		f, err := os.OpenFile(hashFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}

		defer f.Close()

		n, err := d.repo.WriteEmployees(ctx, f)
		if err != nil {
			return "", err
		}

		fmt.Println("written: ", n)

		if err = f.Sync(); err != nil {
			return "", err
		}
	}

	return hashFile, nil
}

func NewDownloaderUC(repo downloader.DWRepository, cfg *config.Config, logger *logrus.Logger) downloader.DwlUC {
	return &DownloaderUC{
		repo:   repo,
		cfg:    cfg.Downloader,
		logger: logger,
	}
}