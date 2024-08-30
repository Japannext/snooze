package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/pkg/apiserver"
	"github.com/japannext/snooze/pkg/processor"
	"github.com/japannext/snooze/pkg/version"
	"github.com/japannext/snooze/pkg/sources/syslog"
	"github.com/japannext/snooze/pkg/sources/otel"
	"github.com/japannext/snooze/pkg/samples"
	"github.com/japannext/snooze/pkg/notifiers/mail"
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
	root.AddCommand(syslog.Cmd)
	root.AddCommand(otel.Cmd)
	root.AddCommand(mail.Cmd)
	root.AddCommand(samples.Cmd)
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
