package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"notification"
	_ "notification/drivers"
	_ "notification/drivers/onesignal"
	_ "notification/drivers/sendgrid"
	_ "notification/drivers/smtp"
	_ "notification/drivers/twilio"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

var Version string
var RunConfig *notification.Config = notification.NewConfig()

const (
	EnvConfigName     = "NOTIFICATION_CONFIG"
	EnvDebugName      = "NOTIFICATION_DEBUG"
	EnvPortName       = "NOTIFICATION_PORT"
	EnvHealthPortName = "NOTIFICATION_HEALTH_PORT"
)

func RunAction(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	cnf := c.String("config")

	dat, err := ioutil.ReadFile(cnf)
	if err != nil {
		return fmt.Errorf("main.go: Config file '%s' read error: %v", cnf, err)
	}

	e := filepath.Ext(cnf)
	if e == ".yml" || e == ".yaml" {
		err = yaml.Unmarshal(dat, RunConfig)
		if err != nil {
			return fmt.Errorf("main.go: Error while unmarshal config file, error: %v", err)
		}
	} else if e == ".json" {
		err = json.Unmarshal(dat, RunConfig)
		if err != nil {
			return fmt.Errorf("main.go: Error while unmarshal config file, error: %v", err)
		}
	} else {
		return errors.New("main.go: Unknown config file format")
	}

	RunConfig.GrpcPort = c.Int("port")
	RunConfig.HealthPort = c.Int("health-port")

	server := notification.NewServer(RunConfig)
	if err := server.LoadDrivers(); err != nil {
		return err
	}
	if err := server.LoadTemplates(); err != nil {
		return err
	}
	return server.ListenAndServe()
}

func main() {
	app := cli.NewApp()
	app.Name = "simple-nofications"
	app.Version = Version
	app.Usage = "Push, Email, SMS notifications with multiple backends"
	app.Author = "Sercan Değirmenci <sercan@otsimo.com>"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config", Value: "config.yml", Usage: "config file path", EnvVar: EnvConfigName},
		cli.IntFlag{Name: "port", Value: notification.DefaultGrpcPort, Usage: "grpc server port", EnvVar: EnvPortName},
		cli.IntFlag{Name: "health-port", Value: notification.DefaultHealthPort, Usage: "health check port", EnvVar: EnvHealthPortName},
		cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log", EnvVar: EnvDebugName},
	}

	app.Action = RunAction
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
