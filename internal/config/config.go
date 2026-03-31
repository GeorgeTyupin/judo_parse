package config

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database   DBConf
	SSH        SSHConf
	Version    string `env:"VERSION"`
	CreateJSON bool   `env:"CREATE_JSON"`
	IsDev      bool   `env:"IS_DEV"`
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

	if err := godotenv.Load("configs/.env", "configs/database.env"); err != nil {
		log.Fatal("Ошибка загрузки .env файлов")
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Ошибка чтения .env файлов")
	}

	return cfg
}

func (d *DBConf) GetConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host,
		d.Port,
		d.Name,
	)
}
