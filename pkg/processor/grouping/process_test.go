package grouping

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/japannext/snooze/pkg/models"
)

func TestProcess(t *testing.T) {
	groupings := []*Grouping{
		{
			If:      `source.Kind == "syslog"`,
			GroupBy: []string{`labels.host`, `labels.process`},
		},
	}
	Startup(groupings)

	tests := []struct {
		Name           string
		Log          *models.Log
		ExpectedLabels map[string]string
		ExpectMatch    bool
	}{
		{
			"syslog hash",
			&api.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog-1.example.com"},
				SeverityText:   "error",
				SeverityNumber: 13,
				Labels:         map[string]string{"host": "host-1", "process": "sshd"},
				Message:        "error: kex_exchange_identification: Connection closed by remote host",
			},
			map[string]string{`labels.host`: "host-1", `labels.process`: "sshd"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := Process(context.TODO(), tt.Log)
			assert.NoError(t, err)
			if tt.ExpectMatch {
				assert.NotEmpty(t, tt.Log.Group.Hash)
			} else {
				assert.Empty(t, tt.Log.Group.Hash)
			}
			for k, v := range tt.ExpectedLabels {
				assert.Equal(t, v, tt.Log.Group.Labels[k])
			}
		})
	}
}
