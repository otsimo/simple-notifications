package notification
import (

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"notification/drivers"
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

func addDriver(dr drivers.Driver) {
	if dr == nil {
		return
	}

}

func LoadDrivers(config *Config) {
	log.Debugf("Config is %v", config)

	for _, r := range (config.Drivers) {
		if driver := drivers.GetDriver(r.Provider); driver != nil {
			if (driver.Type == r.Type) {
				dr, err := driver.New(r.Config)
				if err != nil {
					log.Errorf("Error while creating %s driver:%v", r.Provider, err)
				}else if dr == nil {
					log.Errorf("%s driver is created as nil", r.Provider)
				}else {
					addDriver(dr)
				}
			}else {
				log.Errorf("%s named driver is \"%s\" driver in config, but it is actually a \"%s\" driver", r.Provider, r.Type, driver.Type)
			}
		}else {
			log.Errorf("%s named driver not found", r.Provider)
		}
	}
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