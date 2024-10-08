package googlechat

import (
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "googlechat",
    Short: "Run the snooze-googlechat notification consumer",
    Long:  "Snooze googlechat will consume notifications queues and send googlechat messages",
    Run: func(c *cobra.Command, args []string) {
        Run()
    },
}
