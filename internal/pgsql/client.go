package pgsql

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type Client struct {
	*pgx.Conn
}

func newClient() *Client {

	ctx := context.Background()

	connString := makeConnString()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to pgsql: %s", err)
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
