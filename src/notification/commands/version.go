package commands

import (
	"github.com/spf13/cobra"
	"fmt"
	"path/filepath"
	"os"
	"time"
	"github.com/bugsnag/osext"
	log "github.com/Sirupsen/logrus"

)

var (
	BuildDate = ""
	ApiVersion = "0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Server",
	Long:  `All software has versions. This is Api Server's`,
	Run: func(cmd *cobra.Command, args []string) {
		if BuildDate == "" {
			setBuildDate()
		} else {
			formatBuildDate()
		}
		fmt.Printf("%s\n", ApiVersion)
	},
}

func setBuildDate() {
	fname, _ := osext.Executable()
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		log.Errorln(err)
		return
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		log.Errorln(err)
		return
	}
	t := fi.ModTime()
	BuildDate = t.Format(time.RFC3339)
}

func formatBuildDate() {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", BuildDate)
	BuildDate = t.Format(time.RFC3339)
}