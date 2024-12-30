package lang

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/japannext/snooze/pkg/models"
)

func TestCondition(t *testing.T) {

	tests := []struct {
		Raw         string
		Log         *models.Log
		ExpectMatch bool
	}{
		{
			`source.Kind == "syslog"`,
			&models.Log{Source: models.Source{Kind: "syslog", Name: "prod-syslog"}},
			true,
		},
		{
			`source.Kind == "otlp"`,
			&models.Log{Source: models.Source{Kind: "syslog", Name: "prod-syslog"}},
			false,
		},
		{
			`has(labels.a, labels.b)`,
			&models.Log{Labels: map[string]string{"a": "1", "b": "2"}},
			true,
		},
		{
			`has(labels["c"])`,
			&models.Log{Labels: map[string]string{"a": "1", "b": "2"}},
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
