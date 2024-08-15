package profile

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func TestProcess(t *testing.T) {
	rules := []*Rule{
		{
			If:      `alert.Source.Kind == "syslog"`,
			GroupBy: []string{`alert.Labels.host`, `alert.Labels.process`},
		},
	}
	InitRules(rules)

	tests := []struct {
		Name           string
		Alert          *api.Alert
		ExpectedLabels map[string]string
		ExpectMatch    bool
	}{
		{
			"syslog hash",
			&api.Alert{
				Source:         api.Source{Kind: "syslog", Name: "prod-syslog-1.example.com"},
				SeverityText:   "error",
				SeverityNumber: 13,
				Labels:         map[string]string{"host": "host-1", "process": "sshd"},
				Body:           map[string]string{"body": "error: kex_exchange_identification: Connection closed by remote host"},
			},
			map[string]string{`alert.Labels.host`: "host-1", `alert.Labels.process`: "sshd"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := Process(tt.Alert)
			assert.NoError(t, err)
			if tt.ExpectMatch {
				assert.NotEmpty(t, tt.Alert.GroupHash)
			} else {
				assert.Empty(t, tt.Alert.GroupHash)
			}
			fmt.Printf("GroupHash: %v+\n", tt.Alert.GroupHash)
			fmt.Printf("GroupLabels: %v+\n", tt.Alert.GroupLabels)
			for k, v := range tt.ExpectedLabels {
				assert.Equal(t, v, tt.Alert.GroupLabels[k])
			}
		})
	}
}
