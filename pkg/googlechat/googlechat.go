package googlechat

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/japannext/snooze/pkg/common/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	chat "google.golang.org/api/chat/v1"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	chat *chat.Service
}

func NewClient() *Client {
	ctx := context.Background()

	creds := getCredentials()
	chatService, err := chat.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("failed to initialize chat client: %s", err)
	}

	return &Client{
		chat: chatService,
	}
}

func (c *Client) TestSpaces() {
	ctx := context.Background()
	for name, profile := range profiles {
		req := chat.NewSpacesMessagesService(c.chat).List(profile.Space)
		_, err := req.Context(ctx).Do()
		if err != nil {
			log.Fatalf("failed to load profile '%s' (space='%s'): %s", name, profile.Space, err)
		}
	}
}

func getCredentials() *google.Credentials {
	path := config.ServiceAccountPath
	if path == "" {
		log.Fatal("No SERVICE_ACCOUNT_KEY defined")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading SERVICE_ACCOUNT_KEY file:", err)
	}

	ctx := context.Background()
	cfg, err := google.CredentialsFromJSON(ctx, data, chat.ChatBotScope)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func getHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig:     config.TLS.Config(),
		Proxy:               http.ProxyFromEnvironment,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	// tracing
	tracerProvider := tracing.NewTracerProvider("google")
	transportWithTrace := otelhttp.NewTransport(transport, otelhttp.WithTracerProvider(tracerProvider))

	return &http.Client{
		Transport: transportWithTrace,
	}
}

func (c *Client) SendMessage(ctx context.Context, space string, msg *chat.Message) error {
	req := chat.NewSpacesMessagesService(c.chat).Create(space, msg)
	if _, err := req.Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

var client *Client

func initGooglechat() {
	client = NewClient()
}
