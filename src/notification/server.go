package notification
import (

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
	"fmt"
	"net"
	"google.golang.org/grpc"
)

type Server struct {}

func (s *Server) SendMessage(ctx context.Context, in*Message) (*SendMessageResponse, error) {
	log.Debugln("SendMessage", in.Template)
	fmt.Printf("%d of targets found\n", len(in.Targets))

	for _, t := range (in.Targets) {
		if t.Email != nil {
			fmt.Println("target: Email", t.Email)
		}else if t.Push != nil {
			fmt.Println("target: Push", t.Push)
		}else if t.Sms != nil {
			fmt.Println("target: Sms", t.Sms)
		}
	}
	return &SendMessageResponse{}, nil
}

func LoadDrivers(config *Config) error {
	log.Debugln(config)
	return nil
}

func ListenAndServe(config *Config) {
	port := config.GetPortString()

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterNotificationServiceServer(s, &Server{})

	log.Infoln("Listening", port)

	s.Serve(lis)
}