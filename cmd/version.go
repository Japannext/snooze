package cmd

import (
  "github.com/spf13/cobra"

  "github.com/japannext/snooze/version"
)

var versionCmd = &cobra.Command{
  Use: "version",
  Short: "Display snooze version",
  Long: "Display snooze version",
  Run: func(cmd *cobra.Command, args []string) {
    version.Print()
  },
}
