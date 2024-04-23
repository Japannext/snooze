package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/sources/otel"
	"github.com/japannext/snooze/version"
)

var root = &cobra.Command{
	Use:   "snooze-otel",
	Short: "Alerting system",
	Long:  ``,
}

var cmd = &cobra.Command{
	Use:   "server",
	Short: "Run the snooze-otel server",
	Long:  "Run the snooze-otel server",
	Run: func(c *cobra.Command, args []string) {
		otel.Run()
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
