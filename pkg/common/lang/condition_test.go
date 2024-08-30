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
		Log       *api.Log
		ExpectMatch bool
	}{
		{
			`source.Kind == "syslog"`,
			&api.Log{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			true,
		},
		{
			`source.Kind == "otlp"`,
			&api.Log{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			false,
		},
		{
			`has(labels.a, labels.b)`,
			&api.Log{Labels: map[string]string{"a": "1", "b": "2"}},
			true,
		},
		{
			`has(labels["c"])`,
			&api.Log{Labels: map[string]string{"a": "1", "b": "2"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Raw, func(t *testing.T) {
			c, err := NewCondition(tt.Raw)
			if assert.NoError(t, err) {
				m, err := c.MatchLog(context.Background(), tt.Log)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.ExpectMatch, m)
				}
			}
		})
	}
}
