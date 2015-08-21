package main
import "fmt"
import (
	"text/template"
	"bytes"
	"notification/drivers"
	_ "notification/drivers/sendgrid"
	_ "notification/drivers/mandrill"
	_ "notification/drivers/onesignal"
	_ "notification/drivers/pushwoosh"
	_ "notification/drivers/twilio"
	_ "notification/drivers/mailchimp"
	"google.golang.org/grpc"
	"net"
	"log"
	"notification"
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
	fmt.Println("Hello World!!")

	//templateTest()

	d := drivers.GetDrivers()

	fmt.Println("Drivers:")
	for n, t := range (d) {
		fmt.Printf("%s:%s\n", n, t)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	notification.RegisterNotificationServiceServer(s, &notification.Server{})
	s.Serve(lis)
}