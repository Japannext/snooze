package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/japannext/snooze/pkg/alertmanager"
	"github.com/japannext/snooze/pkg/apiserver"
	"github.com/japannext/snooze/pkg/activecheck"
	"github.com/japannext/snooze/pkg/googlechat"
	"github.com/japannext/snooze/pkg/mail"
	"github.com/japannext/snooze/pkg/otel"
	"github.com/japannext/snooze/pkg/processor"
	"github.com/japannext/snooze/pkg/samples"
	"github.com/japannext/snooze/pkg/syslog"
	"github.com/japannext/snooze/pkg/version"
	"github.com/japannext/snooze/pkg/writer"
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
	root.AddCommand(samples.Cmd)
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
