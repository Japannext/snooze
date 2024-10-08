package googlechat

import (
  "context"
  "net/http"
  "io/ioutil"

  chat "google.golang.org/api/chat/v1"
  log "github.com/sirupsen/logrus"
 "golang.org/x/oauth2"
 "golang.org/x/oauth2/google"
)

func getOauthClient(serviceAccountKeyPath string) *http.Client {
  ctx := context.Background()
  if serviceAccountKeyPath == "" {
    log.Fatal("No SERVICE_ACCOUNT_KEY defined")
  }
  data, err := ioutil.ReadFile(serviceAccountKeyPath)
  if err != nil {
    log.Fatal("Error reading SERVICE_ACCOUNT_KEY file:", err)
  }
  creds, err := google.CredentialsFromJSON(ctx, data, chat.ChatBotScope)
  if err != nil {
    log.Fatal(err)
  }
  return oauth2.NewClient(ctx, creds.TokenSource)
}

type Client struct {
  *chat.SpacesMessagesService
}

func NewClient() *Client {
  c := getOauthClient(config.ServiceAccountPath)
  service, err := chat.New(c)
  if err != nil {
    log.Fatalf("Failed to create chat service: %s", err)
  }
  msgService := chat.NewSpacesMessagesService(service)
  return &Client{SpacesMessagesService: msgService}
}

func (c *Client) SendNewMessage(space, text string) error {
  msg := &chat.Message{
    Text: text,
  }
  if _, err := c.Create(space, msg).Do(); err != nil {
    return err
  }
  return nil
}

func (c *Client) SendReply(space, thread, text string) error {
  msg := &chat.Message{
    Text: text,
  }
  req := c.Create(space, msg).
    MessageReplyOption("REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD").
    ThreadKey(thread)
  if _, err := req.Do(); err != nil {
      return err
  }
  return nil
}

var client *Client

func initGooglechat() {
  client = NewClient()
}
