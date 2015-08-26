package commands

import (
	"notification"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

var RunConfig *notification.Config = notification.NewConfig()
var verbose bool = false

//RootCmd is the root Command
var RootCmd = &cobra.Command{
	Use:   "server",
	Short: "Short Description",
	Long: `Long Description`,
	PersistentPreRun:   func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	},
}

//Execute is runs app
func Execute() {
	addCommands()
	InitializeConfig()
	RootCmd.Execute()
}

func addCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(runCmd)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors:true})

	RootCmd.PersistentFlags().StringP("config", "c", "config.yml", "config file (default is path/config.yml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "debug", "d", false, "enable debug logs")

	//for bash Autocomplete
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	annotation := make(map[string][]string)
	annotation[cobra.BashCompFilenameExt] = validConfigFilenames
	RootCmd.PersistentFlags().Lookup("config").Annotations = annotation
}

func InitializeConfig() {
	viper.SetConfigFile(RootCmd.PersistentFlags().Lookup("config").Value.String())

	err := viper.ReadInConfig()
	viper.AutomaticEnv()

	if err != nil {
		log.Warningln(err, RootCmd.PersistentFlags().Lookup("config").Value.String())
	}

	err = viper.Marshal(RunConfig)
	if err != nil {
		log.Panicln(err)
	}
}