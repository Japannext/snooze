package lang

import (
	"github.com/stretchr/testify/assert"
	"testing"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func TestField(t *testing.T) {
	tests := []struct {
		Raw          string
		Log        *api.Log
		ExpectResult string
	}{
		{
			`source.Kind`,
			&api.Log{Source: api.Source{Kind: "syslog", Name: "prod-syslog"}},
			"syslog",
		},
		{
			`labels["host.name"]`,
			&api.Log{Labels: map[string]string{"host.name": "host-1"}},
			"host-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Raw, func(t *testing.T) {
			f, err := NewField(tt.Raw)
			if assert.NoError(t, err) {
				m, err := ExtractField(tt.Log, f)
				assert.NoError(t, err)
				assert.Equal(t, tt.ExpectResult, m)
			}
		})
	}
}

func TestFieldMap(t *testing.T) {
	tests := []struct{
		name string
		fields []string
		log *api.Log
		expected map[string]string
	}{
		{
			name: "hostproc identity matching",
			fields: []string{"identity.hostname", "identity.process"},
			log: &api.Log{Identity: map[string]string{"hostname": "host01", "process": "sshd"}},
			expected: map[string]string{"identity.hostname": "host01", "identity.process": "sshd"},
		},
		{
			name: "source matching",
			fields: []string{"source.Kind", "source.Name"},
			log: &api.Log{Source: api.Source{Kind: "syslog", Name: "dev"}},
			expected: map[string]string{"source.Kind": "syslog", "source.Name": "dev"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := NewFields(tt.fields)
			assert.NoError(t, err)
			mapper, err := ExtractFields(tt.log, fields)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, mapper)
		})
	}
}
