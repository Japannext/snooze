package lang

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/common/api/v2"
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
	}

	for _, tt := range tests {
		t.Run(tt.Raw, func(t *testing.T) {
			c, err := NewCondition(tt.Raw)
			assert.NoError(t, err)
			m, err := c.Match(context.Background(), tt.Alert)
			assert.NoError(t, err)
			assert.Equal(t, tt.ExpectMatch, m)
		})
	}
}
