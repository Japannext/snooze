package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/japannext/snooze/pkg/common/redis"
	log "github.com/sirupsen/logrus"
)

const SNOOZE_SESSION = "snooze-session"

type Session struct {
	c             *gin.Context  `json:"-"`
	ID            string        `json:"id"`
	Authenticated bool          `json:"authenticated"`
	Username      string        `json:"username"`
	Role          string        `json:"role"`
	Expiration    time.Duration `json:"expiration"`

	// Provider specific values
	OidcState string `json:"oidcState,omitempty"`
}

func sessionKey(id string) string {
	return fmt.Sprintf("sessions/%s", id)
}

// Get the session from redis (if any).
func getSession(c *gin.Context, id string) (*Session, error) {
	ctx := c.Request.Context()
	data, err := redis.Client.Get(ctx, sessionKey(id)).Bytes()
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

// Set the session to redis.
func setSession(ctx context.Context, s *Session) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	if err := redis.Client.Set(ctx, sessionKey(s.ID), data, s.Expiration).Err(); err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}
	return nil
}

func deleteSession(ctx context.Context, id string) error {
	if err := redis.Client.Del(ctx, sessionKey(id)).Err(); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}

// Create a new session.
func newSession(c *gin.Context, ctx context.Context) *Session {
	sessionID := uuid.NewString()
	s := &Session{
		c:          c,
		ID:         sessionID,
		Expiration: time.Duration(6) * time.Hour,
	}
	if err := setSession(ctx, s); err != nil {
		log.Warnf("failed to set session: %s", err)
	}
	c.SetCookie(SNOOZE_SESSION, sessionID, 3600, "/", cookieDomain, true, true)
	return s
}

func (s *Session) Save() error {
	ctx := s.c.Request.Context()
	if err := setSession(ctx, s); err != nil {
		return err
	}
	return nil
}

func (s *Session) Delete(ctx context.Context) {
	if err := deleteSession(ctx, s.ID); err != nil {
		log.Warnf("failed to delete session: %s", err)
	}
	s.c.SetCookie(SNOOZE_SESSION, "", 0, "/", cookieDomain, true, true)
}

func MySession(c *gin.Context, ctx context.Context) *Session {
	sessionID, err := c.Cookie(SNOOZE_SESSION)
	if err != nil { // no session cookie
		log.Warnf("error fetching cookie `%s`: %s", SNOOZE_SESSION, err)
		return newSession(c, ctx)
	}
	s, _ := getSession(c, sessionID)
	if s == nil { // no session in redis
		return newSession(c, ctx)
	}
	return s
}
