package alertcloser

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "alertcloser",
	Short: "Run the snooze-alertcloser job",
	Long:  "Run the snooze-alertcloser job",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
