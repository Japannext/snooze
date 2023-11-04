package server

import (
  "github.com/gocql/gocql"
  gocqlx "github.com/scylladb/gocqlx/v2"
  "github.com/scylladb/gocqlx/v2/table"

  "github.com/japannext/snooze/common/types"
)

type Database struct {
  Cluster *gocql.ClusterConfig
  LogV2 *table.Table
}

// Initialize the ScyllaDB database
func NewDatabase() (*Database, error) {
  hosts := []string{"localhost"}
  cluster := gocql.NewCluster(hosts...)

  db := &Database{
    Cluster: cluster,
    LogV2: table.New(types.LogV2Metadata),
  }

  return db, nil
}

// Return a new session
func (db *Database) NewSession() (gocqlx.Session, error) {
  return gocqlx.WrapSession(db.Cluster.CreateSession())
}
