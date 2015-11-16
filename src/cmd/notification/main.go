package main

import (
	"encoding/json"
	"io/ioutil"
	"notification"
	_ "notification/drivers"
	_ "notification/drivers/mailchimp"
	_ "notification/drivers/mandrill"
	_ "notification/drivers/onesignal"
	_ "notification/drivers/pushwoosh"
	_ "notification/drivers/sendgrid"
	_ "notification/drivers/smtp"
	_ "notification/drivers/twilio"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

var Version string
var RunConfig *notification.Config = notification.NewConfig()

const (
	EnvConfigName = "NOTIFICATION_CONFIG"
	EnvDebugName = "NOTIFICATION_DEBUG"
	EnvPortName = "NOTIFICATION_PORT"
)

func RunAction(c *cli.Context) {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	cnf := c.String("config")

	dat, err := ioutil.ReadFile(cnf)
	if err != nil {
		log.Fatalf("main.go: Config file '%s' read error: %v", cnf, err)
	}

	e := filepath.Ext(cnf)
	if e == ".yml" || e == "yaml" {
		err = yaml.Unmarshal(dat, RunConfig)
		if err != nil {
			log.Fatalf("main.go: Error while unmarshal config file, error: %v", err)
		}
	} else if e == "json" {
		err = json.Unmarshal(dat, RunConfig)
		if err != nil {
			log.Fatalf("main.go: Error while unmarshal config file, error: %v", err)
		}
	} else {
		log.Fatalln("main.go: Unknown config file format")
	}

	envPort := os.Getenv(EnvPortName)

	if len(envPort) > 0 {
		if p, err := strconv.Atoi(envPort); err == nil {
			RunConfig.Port = p
		}
	}

	server := notification.NewServer(RunConfig)
	server.LoadDrivers()
	server.LoadTemplates()
	server.ListenAndServe()
}

func main() {
	app := cli.NewApp()
	app.Name = "simple-nofications"
	app.Version = Version
	app.Usage = "Push, Email, SMS notifications with multiple backends"
	app.Author = "Sercan Değirmenci <sercan@otsimo.com>"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config", Value: "config.yml", Usage: "config file path", EnvVar: EnvConfigName},
		cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log", EnvVar: EnvDebugName},
		cli.IntFlag{Name: "port", Value: notification.DefaultPort, Usage: "grpc server port", EnvVar: EnvPortName},
	}

	app.Action = RunAction
	app.Run(os.Args)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
