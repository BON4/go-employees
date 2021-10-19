package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

func ParseConfig(f string) (Config, error) {
	file, err := os.Open(f)
	if err != nil {

		return Config{}, err
	}
	defer file.Close()

	opts := Config{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(&opts)

	if err != nil {
		return Config{}, err
	}
	return opts, nil
}
