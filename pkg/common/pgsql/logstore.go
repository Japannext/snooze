package pgsql

import (
	"context"
	"os"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/georgysavva/scany/v2/pgxscan"
	log "github.com/sirupsen/logrus"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

var (
	LogStore *PostgresLogStore
)

type PostgresLogStore struct {
	*pgxpool.Pool
}

func NewPostgresLogStore() *PostgresLogStore {
	ctx := context.Background()
	url := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresLogStore{pool}
}

func (pg *PostgresLogStore) Search(query string, pagination api.Pagination) ([]*api.Alert, error) {
	ctx := context.Background()
	var alerts []*api.Alert
	err := pgxscan.Select(ctx, pg, &alerts, `SELECT * FROM snooze_v2_alerts LIMIT 10;`)
	if err != nil {
		return alerts, err
	}
	return alerts, nil
}
