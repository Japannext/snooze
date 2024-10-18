package opensearch

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/japannext/snooze/pkg/common/tracing"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type BootstrapFunction = func(context.Context)

var (
	client *opensearchapi.Client
)

var log *logrus.Entry

func Init() {

	log = logrus.WithFields(logrus.Fields{"module": "opensearch"})

	initConfig()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}
	tracerProvider := tracing.NewTracerProvider("opensearch")
	cfg := opensearch.Config{
		Addresses: config.Addresses,
		Username:  config.Username,
		Password:  config.Password,
		Transport: otelhttp.NewTransport(transport, otelhttp.WithTracerProvider(tracerProvider)),
	}
	var err error
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
