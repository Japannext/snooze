package activecheck

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "activecheck",
	Short: "An active-check (prometheus exporter) daemon for snooze",
	Long:  "An active-check that test several components of snooze",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
