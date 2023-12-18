package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

type postgres struct {
	conn   *pgxpool.Pool
	logger logger.Logger
}

func NewPostgres(host, user, password, db string, logger logger.Logger) IRepository {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   host,
		User:   url.UserPassword(user, password),
		Path:   db,
	}
	conn, err := newConn(dsn)
	if err != nil {
		panic(err)
	}
	return &postgres{
		conn:   conn,
		logger: logger,
	}
}

func newConn(dsn url.URL) (*pgxpool.Pool, error) {
	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()
	conn, err := pgxpool.Connect(context.Background(), dsn.String())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db. %w", err)
	}

	return conn, nil
}
