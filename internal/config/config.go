package config

import (
	_ "embed"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

//go:embed prod.env
var prodEnv string

type Config struct {
	Database DBConf
	SSH      SSHConf
	Version  string `env:"VERSION"`
}

type DBConf struct {
	Host     string `env:"DB_HOST"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Port     int32  `env:"DB_PORT"`
}

type SSHConf struct {
	Host     string `env:"SSH_HOST"`
	User     string `env:"SSH_USER" env-default:"root"`
	Password string `env:"SSH_PASSWORD"`
	Port     string `env:"SSH_PORT" env-default:"22"`
}

func MustLoad() Config {
	var cfg Config

	envMap, err := godotenv.Unmarshal(prodEnv)
	if err != nil {
		log.Fatal("Ошибка чтения конфига")
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Ошибка чтения конфига")
	}

	return cfg
}

func (d *DBConf) GetConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host,
		d.Port,
		d.Name,
	)
}
