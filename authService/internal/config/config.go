package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	GRPC     GrpcConfig    `yaml:"grpc"`
	DB       DBConfig      `yaml:"db"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
}

func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exsist " + path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failes to read config" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
