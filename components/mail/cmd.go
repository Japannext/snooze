package mail

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "mail",
	Short: "Run the snooze-mail notification consumer",
	Long:  "Snooze mail will consume notifications queues and send mails",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
