package tester

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "tester",
	Short: "Populate snooze with sample data",
	Long: "Populate snozoe with sample data for syslog",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
