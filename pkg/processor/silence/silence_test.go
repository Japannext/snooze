package silence_test

import (
	"context"
	"testing"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/silence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	t.Parallel()

	cfg := silence.Config{
		Silences: []silence.Silence{
			{
				Name: "myapp silence",
				If:   `identity.host =~ ".*myapp.*"`,
			},
			{
				Name: "rsyslogd dropper",
				If:   `identity.process == "rsyslogd"`,
				Drop: true,
			},
		},
	}

	p, err := silence.New(cfg)
	require.NoError(t, err)

	tests := []struct{
		Name           string
		Log            *models.Log
		ExpectSkipNotification bool
		ExpectSkipStorage bool
	}{
		{
			"silence match",
			&models.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog.example.com"},
				Identity:       map[string]string{"host": "myapp01", "process": "sshd"},
				Labels:		    map[string]string{},
				SeverityText:   "error",
				SeverityNumber: 13,
				Message:        "error: kex_exchange_identification: Connection closed by remote host",
				Status:			models.LogStatus{},
			},
			true,
			false,
		},
		{
			"silence no match",
			&models.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog.example.com"},
				Identity:       map[string]string{"host": "mysrv01", "process": "sshd"},
				Labels:		    map[string]string{},
				SeverityText:   "error",
				SeverityNumber: 13,
				Message:        "error: kex_exchange_identification: Connection closed by remote host",
				Status:			models.LogStatus{},
			},
			false,
			false,
		},
		{
			"silence drop",
			&models.Log{
				Source:         models.Source{Kind: "syslog", Name: "prod-syslog.example.com"},
				Identity:       map[string]string{"host": "myprogram01", "process": "rsyslogd"},
				Labels:		    map[string]string{},
				SeverityText:   "error",
				SeverityNumber: 13,
				Message:        "this is an example",
				Status:			models.LogStatus{},
			},
			true,
			true,
		},
	}

    for _, tt := range tests {
        t.Run(tt.Name, func(t *testing.T) {
            t.Parallel()

            p.Process(context.TODO(), tt.Log)

            assert.Equal(t, tt.ExpectSkipNotification, tt.Log.Status.SkipNotification)
            assert.Equal(t, tt.ExpectSkipStorage, tt.Log.Status.SkipStorage)
        })
    }
}
