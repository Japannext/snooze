package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/processor"
	"github.com/japannext/snooze/version"
)

var root = &cobra.Command{
	Use:   "snooze-processor",
	Short: "Process alerts along trnasform pipelines",
	Long:  ``,
}

var cmd = &cobra.Command{
	Use:   "server",
	Short: "Run the snooze-processor service",
	Long:  "Run the snooze-processor service",
	Run: func(c *cobra.Command, args []string) {
		processor.Run()
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
