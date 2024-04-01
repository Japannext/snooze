package version

import (
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
  Use: "version",
  Short: "Display snooze version",
  Long: "Display snooze version",
  Run: func(cmd *cobra.Command, args []string) {
    Print()
  },
}
