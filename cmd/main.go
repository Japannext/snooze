package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/components/activecheck"
	"github.com/japannext/snooze/components/alertmanager"
	"github.com/japannext/snooze/components/apiserver"
	"github.com/japannext/snooze/components/googlechat"
	"github.com/japannext/snooze/components/mail"
	"github.com/japannext/snooze/components/otel"
	"github.com/japannext/snooze/components/processor"
	"github.com/japannext/snooze/components/syslog"
	"github.com/japannext/snooze/components/writer"
	"github.com/japannext/snooze/internal/version"
)

var root = &cobra.Command{
	Use:   "snooze",
	Short: "",
	Long:  ``,
}

func init() {
	root.AddCommand(alertmanager.Cmd)
	root.AddCommand(apiserver.Cmd)
	root.AddCommand(activecheck.Cmd)
	root.AddCommand(googlechat.Cmd)
	root.AddCommand(mail.Cmd)
	root.AddCommand(otel.Cmd)
	root.AddCommand(processor.Cmd)
	root.AddCommand(syslog.Cmd)
	root.AddCommand(version.Cmd)
	root.AddCommand(writer.Cmd)
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
