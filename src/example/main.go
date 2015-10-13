package main

import (
	"log"

	pb "notificationpb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:18844"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotificationServiceClient(conn)

	message := &pb.Message{
		Template: "welcome",
		Language: "tr",
		Targets: []*pb.Target{
			pb.NewEmailTarget(&pb.Email{
				ToEmail:   []string{"degirmencisercan@gmail.com"},
				FromEmail: "noreply@otsimo.com",
				FromName:  "Sercan Değirmenci",
				Subject:   "Hello World!!",
			}),
			pb.NewSmsTarget(&pb.Sms{
				To: []string{"+21123124", "+123124"},
			}),
			pb.NewPushTarget(&pb.Push{
				To: []string{"asdaf78a6sfa6f5asf", "j1g24feqfwd7as6d6t7asf"},
			}),
		},
	}

	r, err := c.SendMessage(context.Background(), message)

	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	log.Printf("Result: %d \n %s \n", r.Type, r.Data)
}
