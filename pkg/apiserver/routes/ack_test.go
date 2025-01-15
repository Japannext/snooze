package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	opensearch.Init()
	redis.Init()
	mq.Init()
	os.Exit(m.Run())
}

func TestPostAck(t *testing.T) {
	router := gin.Default()
	router.POST("/api/ack", postAck)
	w := httptest.NewRecorder()

	data := `{
		"time": 1736983672000,
		"username": "john.doe",
		"reason": "noisy alert",
		"logIDs": ["abc123"]
	}`
	req, _ := http.NewRequest("POST", "/api/ack", strings.NewReader(data))
	router.ServeHTTP(w, req)

	opensearch.Flush(context.Background(), models.AckIndex)

	assert.Equal(t, 201, w.Code)
}
