package processor

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "processor",
	Short: "Run the snooze-processor service",
	Long:  "Run the snooze-processor service",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
