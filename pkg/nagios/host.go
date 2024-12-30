package nagios

import (
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	icinga   bool
	hostname string
	state    string
	labels   map[string]string
	ts       string
)

const TIME_FORMAT = "2006-01-02 15:04:05 +0000"

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Send a host notification to snooze",
	Run:   sendHost,
}

func init() {
	flagset := hostCmd.PersistentFlags()
	flagset.StringVarP(&hostname, "hostname", "H", "HOSTNAME", "The host that emit the alert")
	flagset.StringVarP(&state, "state", "S", "STATE", "The nagios state the alert is in")
	flagset.BoolVarP(&icinga, "icinga", "i", false, "Whether it's icinga or nagios")

	flagset.StringVarP(&ts, "timestamp", "t", "", "The timestamp of the event when it was observed")

	flagset.StringToStringVarP(&labels, "labels", "l", map[string]string{}, "Custom labels to be added to the alert")

	// Mandatory flags
	hostCmd.MarkFlagRequired("hostname")
	hostCmd.MarkFlagRequired("state")

	Cmd.AddCommand(hostCmd)
}

func sendHost(c *cobra.Command, args []string) {
	item := &models.Alert{}

	if icinga {
		item.Source.Kind = "icinga"
		item.AlertGroup = "IcingaHost"
	} else {
		item.Source.Kind = "nagios"
		item.AlertGroup = "NagiosHost"
	}

	item.AlertName = "Host"

	item.Summary = fmt.Sprintf("Host if %s", state)

	if ts != "" {
		t, err := time.Parse(TIME_FORMAT, ts)
		if err != nil {
			log.Fatalf("failed to parse timestamp '%s': %s", t, err)
		}
		item.StartAt = models.Time{Time: t}
	}

	for key, value := range labels {
		item.Labels[key] = value
	}

	err := sendAlert(item)
	if err != nil {
		log.Fatal(err)
	}
}
