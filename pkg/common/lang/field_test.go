package lang

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func TestField(t *testing.T) {
	tests := []struct {
		Raw          string
		Alert        *api.Alert
		ExpectResult string
	}{
		{
			`alert.Source.Kind`,
			&api.Alert{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			"syslog",
		},
		{
			`alert.Labels["host.name"]`,
			&api.Alert{Labels: map[string]string{"host.name": "host-1"}},
			"host-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Raw, func(t *testing.T) {
			f, err := NewField(tt.Raw)
			if assert.NoError(t, err) {
				m, err := f.Get(context.Background(), tt.Alert)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectResult, m)
			}
		})
	}
}
