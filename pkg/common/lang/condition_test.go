package lang

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func TestCondition(t *testing.T) {

	tests := []struct {
		Raw         string
		Alert       *api.Alert
		ExpectMatch bool
	}{
		{
			`alert.Source.Kind == "syslog"`,
			&api.Alert{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			true,
		},
		{
			`alert.Source.Kind == "otlp"`,
			&api.Alert{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			false,
		},
		{
			`has(alert.Labels.a, alert.Labels.b)`,
			&api.Alert{Labels: map[string]string{"a": "1", "b": "2"}},
			true,
		},
		{
			`has(alert.Labels["c"])`,
			&api.Alert{Labels: map[string]string{"a": "1", "b": "2"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Raw, func(t *testing.T) {
			c, err := NewCondition(tt.Raw)
			if assert.NoError(t, err) {
				m, err := c.Match(context.Background(), tt.Alert)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.ExpectMatch, m)
				}
			}
		})
	}
}
