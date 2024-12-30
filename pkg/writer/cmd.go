package writer

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "writer",
	Short: "Write items in queue to opensearch database",
	Long:  "Interface between a nats queue and opensearch to batch insert items.",
	Run: func(c *cobra.Command, args []string) {
		Run()
	},
}
