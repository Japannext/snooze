package exporter

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "exporter",
	Short: "An active-check exporter for snooze",
	Long: "An active-check exporter that test several components of snooze",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
