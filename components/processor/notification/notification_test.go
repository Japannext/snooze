package notification_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/components/processor/notification"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/redis"
)

func TestMain(m *testing.M) {
	opensearch.Init()
	redis.Init()
	mq.Init()
	os.Exit(m.Run())
}

func TestProcess(t *testing.T) {
	t.Parallel()

	cfgs := []pipeline.Notification{
		{
			Name: "sre-team",
			If:   `labels.owner == "sre"`,
			Destinations: []models.Destination{
				{Queue: "googlechat", Profile: "sre"},
			},
		},
	}

	defaultDests := []models.Destination{
		{Queue: "mail", Profile: "sysadmin"},
	}

	p, err := notification.New(cfgs, defaultDests)
	require.NoError(t, err)

	tests := []struct {
		Name                 string
		Log                  *models.Log
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
