package snooze_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/sirupsen/logrus/hooks/test"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/logging"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/processor/snooze"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logging.Init()
	mq.Init()
	opensearch.Init()
	redis.Init()
	os.Exit(m.Run())
}

func newSnoozes(items []models.Snooze) *redis.TestKeys {
	kvs := map[string]string{}
	for _, item := range items {
		value, _ := json.Marshal(item)
		for _, group := range item.Groups {
			if group.Hash == "" {
				group.Hash = utils.ComputeHash(group.Labels)
			}
			key := fmt.Sprintf("snooze/%s/%s", group.Name, group.Hash)
			kvs[key] = string(value)
		}
	}

	return redis.NewKeys(kvs)
}

func newGroup(name string, labels map[string]string) *models.Group {
	return &models.Group{
		Name: name,
		Labels: labels,
		Hash: utils.ComputeHash(labels),
	}
}

func TestProcess(t *testing.T) {
	p := snooze.Processor{}

	startAt := models.Time{Time: time.Now()}
	endAt := models.Time{Time: time.Now().Add(time.Hour)}

	testkeys := newSnoozes([]models.Snooze{
		{
			Groups: []models.Group{
				{Name: "by-host", Labels: map[string]string{"identity.hostname": "myapp01"}},
				{Name: "by-host", Labels: map[string]string{"identity.hostname": "myapp02"}},
			},
			Reason: "myapp* servers are noisy",
			StartsAt: startAt,
			EndsAt: endAt,
			Username: "john.doe",
		},
	})
	defer testkeys.Cleanup()

	tests := []struct{
		Name string
		Log *models.Log
		ExpectSkipNotification bool
	}{
		{
			Name: "match",
			Log: &models.Log{
				Groups: []*models.Group{
					newGroup("by-host", map[string]string{"identity.hostname": "myapp01"}),
				},
			},
			ExpectSkipNotification: true,
		},
		{
			Name: "no match",
			Log: &models.Log{
				Groups: []*models.Group{
					newGroup("by-host", map[string]string{"identity.hostname": "myapp03"}),
				},
			},
			ExpectSkipNotification: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := p.Process(context.TODO(), tt.Log)
			assert.Equal(t, decision.KindOK, res.Kind)
			assert.Equal(t, tt.ExpectSkipNotification, tt.Log.Status.SkipNotification)
		})
	}
}
