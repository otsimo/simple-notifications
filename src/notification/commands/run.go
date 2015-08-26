package commands
import (
	"github.com/spf13/cobra"
	"notification"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Long:  `Run the Server with given config`,
	Run:   func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		notification.LoadDrivers(RunConfig)
		notification.ListenAndServe(RunConfig)
	},
}

