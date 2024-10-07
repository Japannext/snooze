package pgsql

import (
	"context"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/georgysavva/scany/v2/pgxscan"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
)

var (
	LogStore *PostgresLogStore
)

type PostgresLogStore struct {
	pool *pgxpool.Pool
	conn *pgxpool.Conn
}

func NewPostgresLogStore() *PostgresLogStore {
	ctx := context.Background()
	cfg := initConfig()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgxpool.Acquire(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresLogStore{pool, conn}
}

func (pg *PostgresLogStore) Get(uid string) (*models.Log, error) {
	ctx := context.Background()

	rows, err := pg.conn.Query(ctx, fmt.Sprintf(`SELECT * FROM log_v2 WHERE uid = '%s' LIMIT 1`, uid))
	if err != nil {
		return nil, err
	}
	var item *models.Log
	err := pgxscan.ScanOne(&item, rows)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (pg *PostgresLogStore) Search(query string, pagination models.Pagination) ([]*models.Log, error) {
	ctx := context.Background()
	var aitems []*models.Log
	err := pgxscan.Select(ctx, pg, &alerts, `SELECT * FROM snooze_v2_alerts LIMIT 10;`)
	if err != nil {
		return alerts, err
	}
	return alerts, nil
}
