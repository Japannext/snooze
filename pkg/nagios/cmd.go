package nagios

import (
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nagios",
	Short: "Send a nagios/icinga notification to snooze",
	Long:  "Send a nagios/icinga notification to snooze",
}

var (
	snoozeURL *url.URL
)

func init() {
	var u string
	Cmd.PersistentFlags().StringVarP(&u, "url", "u", "", "URL to the snooze endpoint")
	Cmd.MarkFlagRequired("url")

	var err error
	snoozeURL, err = url.Parse(u)
	if err != nil {
		log.Fatalf("failed to parse URL '%s': %s", u, err)
	}
}
