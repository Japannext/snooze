package grouping_test

import (
	"context"
	"os"
	"fmt"
	"testing"

	"github.com/japannext/snooze/pkg/processor/grouping"
	"github.com/japannext/snooze/pkg/processor/decision"
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/utils"
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

func findGroup(item *models.Log, name string) (bool, *models.Group) {
	for _, gr := range item.Groups {
		if gr.Name == name {
			return true, gr
		}
	}

	return false, &models.Group{}
}

func groupKey(group *models.Group) string {
	hash := utils.ComputeHash(group.Labels)

	return fmt.Sprintf("grouping:%s:%s", group.Name, hash)
}

func TestProcess(t *testing.T) {
	t.Parallel()

	cfg := grouping.Config{
		Groupings: []grouping.Grouping{
			{
				Name:    "by-host-process",
				If:      `source.kind == "syslog"`,
				GroupBy: []string{`labels.host`, `labels.process`},
			},
		},
	}

	p, err := grouping.New(cfg)
	require.NoError(t, err)

	tests := []struct {
		Name            string
		Log             *models.Log
		ExpectedLabels  map[string]string
		ExpectedMatches []string
	}{
		{
			"syslog hash",
			&models.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog-1.example.com"},
				SeverityText:   "error",
				SeverityNumber: 13,
				Labels:         map[string]string{"host": "host-1", "process": "sshd"},
				Message:        "error: kex_exchange_identification: Connection closed by remote host",
			},
			map[string]string{`labels.host`: "host-1", `labels.process`: "sshd"},
			[]string{"by-host-process"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			t.Cleanup(func() {
				keys := []string{}
				for _, group := range tt.Log.Groups {
					keys = append(keys, groupKey(group))
				}

				redis.Client.Del(context.TODO(), keys...)
			})

			res := p.Process(context.TODO(), tt.Log)
			assert.Equal(t, decision.KindOK, res.Kind)

			for _, name := range tt.ExpectedMatches {
				hasGroup, group := findGroup(tt.Log, name)
				assert.True(t, hasGroup)

				if hasGroup {
					for k, v := range tt.ExpectedLabels {
						assert.Equal(t, v, group.Labels[k])
					}
				}
			}
		})
	}
}
