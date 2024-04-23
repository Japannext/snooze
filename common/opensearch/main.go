package opensearch

import (
	v2 "github.com/opensearch-project/opensearch-go/v2"
	"github.com/sirupsen/logrus"
)

var (
	Client *OpensearchClient
)

type OpensearchClient struct {
	*v2.Client
}

var log *logrus.Entry

func Init() {

	log = logrus.WithFields(logrus.Fields{"module": "opensearch"})

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
