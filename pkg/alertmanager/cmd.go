package alertmanager

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "alertmanager",
	Short: "Run the snooze-alertmanager source",
	Long:  "Run the snooze-alertmanager source",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
