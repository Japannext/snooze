package syslog

import (
    "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
    Use:   "syslog",
    Short: "Run the snooze-syslog service",
    Long:  "Run the snooze-syslog service",
    Run: func(c *cobra.Command, args []string) {
        Run()
    },
}
