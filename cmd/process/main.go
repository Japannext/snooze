package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/japannext/snooze/process"
  "github.com/japannext/snooze/cmd"
)

var rootCmd = &cobra.Command{
  Use: "snooze-process",
  Short: "",
  Long: ``,
}

var processCmd = &cobra.Command{
  Use: "process",
  Short: "Run the snooze-process server",
  Long: "Run the snooze-process server",
  Run: func(cmd *cobra.Command, args []string) {
    server.Run()
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.AddCommand(processCmd)
  rootCmd.AddCommand(cmd.versionCmd)
}
