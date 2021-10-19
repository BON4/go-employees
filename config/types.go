package config

import "time"

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host           string `yaml:"host"`
	Port           uint16 `yaml:"port"`
	Database       string `yaml:"database"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	SSLMode		   string `yaml:"ssl_mode"`
	ConnectTimeout time.Duration `yaml:"connect-timeout"`
}

