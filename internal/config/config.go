package config

import (
	"fmt"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database   DBConf
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

func MustLoad() *Config {
	var cfg Config

	godotenv.Load("configs/.env", "configs/database.env")
	cleanenv.ReadEnv(&cfg)

	return &cfg
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
