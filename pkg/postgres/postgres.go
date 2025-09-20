package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/projectWithOrder/configs"
)

func NewPostgresDB(cfg configs.Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s database=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
