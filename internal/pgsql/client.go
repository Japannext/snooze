package pgsql

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/tern/v2/migrate"
	_ "github.com/lib/pq"
	"github.com/wcamarao/pmx"
)

type Client struct {
	*pgx.Conn
}

//go:embed migrations/*.sql
var migrationFS embed.FS

func newClient() *Client {

	ctx := context.Background()

	connString := makeConnString()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to pgsql: %s", err)
	}

	migrator, err := migrate.NewMigrator(context.Background(), conn, "my_schema_version")
	if err != nil {
		log.Fatalf("failed to initialize migrator: %s", err)
	}

	if err := migrator.LoadMigrations(migrationFS); err != nil {
		log.Fatalf("failed to load migrations: %s", err)
	}

	if err := migrator.Migrate(ctx); err != nil {
		log.Fatalf("failed to execute migration: %s", err)
	}

	return &Client{
		Conn: conn,
	}
}

func makeConnString() string {
	config := GetConfig()
	return fmt.Sprintf("host='%s' user='%s' port=%d database='%s'", config.Host, config.Username, config.Port, config.Database)
}

var client *Client
var clientOnce sync.Once

func GetClient() *Client {
	clientOnce.Do(func() {
		client = newClient()
	})

	return client
}

// Insert an object into postgreSQL
func Insert(ctx context.Context, item any) (pgconn.CommandTag, error) {
	return pmx.Insert(ctx, GetClient().Conn, item)
}

func Select(ctx context.Context, dest any, sql string) error {
	return pmx.Select(ctx, GetClient().Conn, dest, sql)
}
