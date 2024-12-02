package sessions

import (
	"context"
	"fmt"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/redis"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("snooze")
}

type Session struct {
	c *gin.Context `json:"-"`
	ID string `json:"id"`
	Authenticated bool `json:"authenticated"`
	Username string `json:"username"`
	Role string `json:"role"`
	Expiration time.Duration `json:"expiration"`
	AuthProvider string `json:"authProvider,omitempty"`
	AuthSession json.RawMessage `json:"authSession,omitempty"`
}

// Get the session from redis (if any)
func getSession(c *gin.Context, id string) (*Session, error) {
	ctx := c.Request.Context()
	key := fmt.Sprintf("session/%s", id)
	data, err := redis.Client.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	var s *Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session from redis: %w", err)
	}
	s.c = c
	return s, nil
}

// Set the session to redis
func setSession(ctx context.Context, s *Session) error {
	key := fmt.Sprintf("sessions/%s", s.ID)
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	if err := redis.Client.Set(ctx, key, data, s.Expiration).Err(); err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}
	return nil
}

func deleteSession(ctx context.Context, id string) error {
	key := fmt.Sprintf("sessions/%s", id)
	if err := redis.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}

// Create a new session
func newSession(c *gin.Context, ctx context.Context) *Session {
	sessionID := uuid.NewString()
	s := &Session{
		c: c,
		ID: sessionID,
		Expiration: time.Duration(6) * time.Hour,
	}
	if err := setSession(ctx, s); err != nil {
		log.Warnf("failed to set session: %s", err)
	}
	c.SetCookie("snooze-session", sessionID, 3600, "/", "localhost", true, true)
	return s
}

func (s *Session) Save() error {
	ctx := s.c.Request.Context()
	if err := setSession(ctx, s); err != nil {
		return err
	}
	return nil
}

func (s *Session) Delete() {
	ctx := s.c.Request.Context()
	if err := deleteSession(ctx, s.ID); err != nil {
		log.Warnf("failed to delete session: %s", err)
	}
	s.c.SetCookie("snooze-session", "", 0, "/", "localhost", true, true)
}

func MySession(c *gin.Context) *Session {
	ctx, span := tracer.Start(c.Request.Context(), "MySession")
	defer span.End()

	sessionID, err := c.Cookie("snooze-session")
	if err != nil { // no session cookie
		return newSession(c, ctx)
	}
	s, _ := getSession(c, sessionID)
	if s == nil { // no session in redis
		return newSession(c, ctx)
	}
	return s
}
