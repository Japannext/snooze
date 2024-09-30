package opensearch

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type BootstrapFunction = func(context.Context)

var (
	client *opensearchapi.Client
)

var log *logrus.Entry

func Init() {

	log = logrus.WithFields(logrus.Fields{"module": "opensearch"})

	// client initialization
	cfg, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err = opensearchapi.NewClient(opensearchapi.Config{Client: cfg})
	if err != nil {
		log.Fatal(err)
	}

	// Fail immediately if the database is unreachable
	ctx := context.Background()
	if err := CheckHealth(ctx); err != nil {
		log.Fatal(err)
	}

	// bootstrap
	bootstrap(ctx)
}
