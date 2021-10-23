package config

import "time"

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Downloader Downloader `yaml:"downloader"`
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	AppVersion        string        `yaml:"app_version"`
	Port              string        `yaml:"port"`
	PprofPort         string        `yaml:"pprof_port"`
	Mode              string        `yaml:"mode"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	SSL               bool          `yaml:"ssl"`
	CtxDefaultTimeout time.Duration `yaml:"ctx_default_timeout"`
	Debug             bool          `yaml:"debug"`
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

type Downloader struct {
	FileFolder string `yaml:"file_folder"`
}
