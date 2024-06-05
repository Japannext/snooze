package apiserver

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Run the snooze-apiserver",
	Long:  "Run the snooze-apiserver",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
