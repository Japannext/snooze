package notification_test

import (
	"context"
	"os"
	"testing"

	"github.com/japannext/snooze/pkg/processor/notification"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	opensearch.Init()
	redis.Init()
	mq.Init()
	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	t.Parallel()

	cfg := notification.Config{
		Notifications: []notification.Notification{
			{
				Name:         "sre-team",
				If:           `labels.owner == "sre"`,
				Destinations: []models.Destination{
					{Queue: "googlchat", Profile: "sre"},
				},
			},
		},
	}

	p, err := notification.New(cfg)
	require.NoError(t, err)

	tests := []struct {
		Name                string
		Log                 *models.Log
		ExpectedDestinations []models.Destination
	}{
		{
			"sre alert",
			&models.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog-1.example.com"},
				SeverityText:   "error",
				SeverityNumber: 13,
				Labels:         map[string]string{"owner": "sre"},
				Message:        "error: kex_exchange_identification: Connection closed by remote host",
			},
			[]models.Destination{
				{Queue: "googlechat", Profile: "sre"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			res := p.Process(context.TODO(), tt.Log)
			assert.Equal(t, decision.KindOK, res.Kind)
		})
	}
}
