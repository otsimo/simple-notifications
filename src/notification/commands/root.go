package commands

import (
	"notification"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

var RunConfig *notification.Config = notification.NewConfig()

//RootCmd is the root Command
var RootCmd = &cobra.Command{
	Use:   "server",
	Short: "Short Description",
	Long: `Long Description`,
}

//Execute is runs app
func Execute() {
	addCommands()
	RootCmd.Execute()
}

func addCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(runCmd)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, ForceColors:true})

	RootCmd.PersistentFlags().StringP("config", "c", "config.yml", "config file (default is path/config.yml)")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	//for bash Autocomplete
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	annotation := make(map[string][]string)
	annotation[cobra.BashCompFilenameExt] = validConfigFilenames
	RootCmd.PersistentFlags().Lookup("config").Annotations = annotation
}

func InitializeConfig() {
	viper.SetConfigFile(RootCmd.PersistentFlags().Lookup("config").Value.String())
	//viper.AddConfigPath(Source)
	err := viper.ReadInConfig()
	if err != nil {
		log.Warningln(err, RootCmd.PersistentFlags().Lookup("config").Value.String())
	}
	err = viper.Marshal(RunConfig)
	if err != nil {
		log.Panicln(err)
	}
}