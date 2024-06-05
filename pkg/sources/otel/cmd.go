package otel

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "otel",
	Short: "Run the snooze-otel server",
	Long:  "Run the snooze-otel server",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
