package pkg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
	Sslmode  string
}

func InitPostgres(cfg Config) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("host =%s port =%s user =%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Dbname, cfg.Password, cfg.Sslmode))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
