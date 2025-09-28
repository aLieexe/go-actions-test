package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(DBURL string) (*pgxpool.Pool, error) {
	stringConnection := DBURL
	config, err := pgxpool.ParseConfig(stringConnection)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
