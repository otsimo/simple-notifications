package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"notification"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := notification.NewNotificationServiceClient(conn)

	message := &notification.Message{
		Template: "welcome",
		Language: "tr",
		Targets: []*notification.Target{
			&notification.Target{
				Email: &notification.Email{
					ToEmail:   []string{"to1@examp.com", "to2@examp.com"},
					FromEmail: "from@asd.com",
				},
			},
			&notification.Target{
				Sms: &notification.Sms{
					To: []string{"+21123124", "+123124"},
				},
			},
			&notification.Target{
				Push: &notification.Push{
					To: []string{"asdaf78a6sfa6f5asf", "j1g24feqfwd7as6d6t7asf"},
				},
			},
		},
	}

	r, err := c.SendMessage(context.Background(), message)

	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	log.Printf("Result: %d \n %s \n", r.Type, r.Data)
}
