package main

import "fmt"
import (
	"bytes"
	_ "google.golang.org/grpc"
	_ "net"
	_ "notification"
	"notification/commands"
	_ "notification/drivers"
	_ "notification/drivers/mailchimp"
	_ "notification/drivers/mandrill"
	_ "notification/drivers/onesignal"
	_ "notification/drivers/pushwoosh"
	_ "notification/drivers/sendgrid"
	_ "notification/drivers/twilio"
	"text/template"
)

func templateTest() {
	const letter = `Dear {{.name}}, And {{.gift}},
	{{.attended}}
	Best wishes,
	SD
	`
	t := template.Must(template.New("letter").Parse(letter))

	commits := map[string]interface{}{
		"name":     "Sercan Degirmenci",
		"attended": true,
		"gift":     "xyaz",
	}

	var doc bytes.Buffer
	err := t.Execute(&doc, commits)
	if err == nil {
		s := doc.String()
		fmt.Println(s)
	} else {
		panic(err)
	}
}

var Version string

func main() {
	commands.ApiVersion = Version
	commands.Execute()
}
