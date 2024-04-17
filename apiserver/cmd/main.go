package main

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/japannext/snooze/apiserver"
  "github.com/japannext/snooze/version"
)

var root = &cobra.Command{
  Use: "snooze-apiserver",
  Short: "Alerting system",
  Long: ``,
}

var cmd = &cobra.Command{
  Use: "server",
  Short: "Run the snooze-api server",
  Long: "Run the snooze-api server",
  Run: func(c *cobra.Command, args []string) {
    apiserver.Run()
  },
}

func init() {
  root.AddCommand(cmd)
  root.AddCommand(version.Cmd)
}

func main() {
  if err := root.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
