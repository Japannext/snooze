package silence

import (
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/common/api/v2"
)

func TestProcess(t *testing.T) {
	rules := []*Rule{
		{
			Name: "syslog",
			If: `alert.Source.Kind == "syslog"`,
		},
	}
	InitRules(rules)

	tests := []struct {
		Name           string
		Alert          *api.Alert
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
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := Process(tt.Alert)
			if assert.NoError(t, err) {
				if tt.ExpectMatch {
					assert.Equal(t, true, tt.Alert.Mute.Enabled)
					assert.Equal(t, "silence", tt.Alert.Mute.Component)
					assert.Equal(t, "syslog", tt.Alert.Mute.Rule)
					assert.Equal(t, true, tt.Alert.Mute.SkipNotification)
				} else {
					assert.Equal(t, false, tt.Alert.Mute.Enabled)
				}
			}
		})
	}
}
