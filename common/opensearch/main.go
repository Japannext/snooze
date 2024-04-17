package opensearch

import (
  log "github.com/sirupsen/logrus"
  v2 "github.com/opensearch-project/opensearch-go/v2"
)

var (
  Client *OpensearchClient
)

type OpensearchClient struct {
  *v2.Client
}

func Init(check bool) {

  cfg, err := initConfig()
  if err != nil {
    log.Fatal(err)
  }

  client, err := v2.NewClient(*cfg)
  if err != nil {
    log.Fatal(err)
  }

  Client = &OpensearchClient{client}

  // Fail immediately if the database is unreachable
  if err := Client.CheckHealth(); err != nil {
    log.Fatal(err)
  }

  if err := Client.Bootstrap(); err != nil {
    log.Fatal(err)
  }
}
