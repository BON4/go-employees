package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

func ParseEnv(c Config) (Config, error) {
	c.Downloader.FileFolder = os.Getenv("CSV_FLD")
	if len(c.Downloader.FileFolder) == 0 {
		return Config{}, errors.New("no folder for tables in csv is provided")
	}
	return c, nil
}

func ParseConfig(f string) (Config, error) {
	file, err := os.Open(f)
	if err != nil {

		return Config{}, err
	}
	defer file.Close()

	opts := Config{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(&opts)

	opts, err = ParseEnv(opts)

	if err != nil {
		return Config{}, err
	}
	return opts, nil
}
