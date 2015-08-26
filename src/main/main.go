package main
import "fmt"
import (
	"text/template"
	"bytes"
	_ "notification/drivers"
	_ "notification/drivers/sendgrid"
	_ "notification/drivers/mandrill"
	_ "notification/drivers/onesignal"
	_ "notification/drivers/pushwoosh"
	_ "notification/drivers/twilio"
	_ "notification/drivers/mailchimp"
	_ "google.golang.org/grpc"
	_ "net"
	_ "notification"
	"notification/commands"
)

func templateTest() {
	const letter =
	`Dear {{.name}}, And {{.gift}},
	{{.attended}}
	Best wishes,
	SD
	`
	t := template.Must(template.New("letter").Parse(letter))

	commits := map[string]interface{}{
		"name": "Sercan Degirmenci",
		"attended":   true,
		"gift": "xyaz",
	}

	var doc bytes.Buffer
	err := t.Execute(&doc, commits)
	if err == nil {
		s := doc.String()
		fmt.Println(s)
	}else {
		panic(err)
	}

}
func main() {
	commands.Execute()
}