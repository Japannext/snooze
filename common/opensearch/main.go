package opensearch

import (
  v2 "github.com/opensearch-project/opensearch-go/v2"
)

type Database struct {
  Client *v2.Client
}

func Init() (*Database, error) {

  cfg, err := initConfig()
  if err != nil {
    return &Database{}, err
  }

  client, err := v2.NewClient(cfg.v2Config())
  if err != nil {
    return &Database{}, err
  }

  db := &Database{client}

  if err := db.Bootstrap(); err != nil {
    return db, err
  }

  return db, nil
}
