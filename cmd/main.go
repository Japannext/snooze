package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/apiserver"
	"github.com/japannext/snooze/processor"
	"github.com/japannext/snooze/version"
)

var root = &cobra.Command{
	Use:   "snooze",
	Short: "",
	Long:  ``,
}

func init() {
	root.AddCommand(version.Cmd)
	root.AddCommand(processor.Cmd)
	root.AddCommand(apiserver.Cmd)
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
