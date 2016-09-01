package main

import (
	"log"

	pb "notificationpb"

	"time"

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
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	sr, _ := c.Scan(context.Background(), &pb.ScanRequest{})
	log.Printf("Scan %v:", sr.Events)
	message := &pb.Message{
		Event:    "welcome",
		Language: "en",
		Tags:     map[string]string{"user_id": "b5980760-dc0b-4ed9-9aa5-489c85c5fa5e"},
		DataJson: pb.Map2Str(map[string]interface{}{
			"name":  "Sercan",
			"count": 1}),
		Targets: pb.NewTargets(
			&pb.Email{
				ToEmail: []string{"degirmencisercan@gmail.com"},
				Cc:      []string{"sercan@otsimo.com"},
			},
			&pb.Push{
				Template: []byte(`Selam fella`),
			}, &pb.Sms{}),
	}
	r, err := c.SendMessage(ctx, message)
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	log.Printf("Result: %s\n%d", r.Output, len(r.Results))
	for _, r2 := range r.Results {
		log.Printf("[%s]: %s", r2.Target, r2.Output)
	}
}
