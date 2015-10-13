package notification

import (
	"net"
	"notification/drivers"
	pb "notificationpb"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	Config  *Config
	Drivers map[string]drivers.Driver
}

func (server *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.SendMessageResponse, error) {
	log.Debugln("SendMessage", in.Template)

	for _, t := range in.Targets {
		ty := getTargetType(t)
		if len(ty) > 0 {
			driver := server.Drivers[ty]
			if driver != nil {
				driver.Send(in, t)
			} else {
				log.Warningf("%s driver not found", ty)
			}
		} else {
			log.Warningf("there is no suitable driver")
		}
	}
	return &pb.SendMessageResponse{}, nil
}

func (server *Server) addDriver(dr drivers.Driver) {
	server.Drivers[dr.Type()] = dr
}

func (server *Server) LoadDrivers() {
	log.Debugf("Config is %v", server.Config)

	for _, r := range server.Config.Drivers {
		if driver := drivers.GetDriver(r.Provider); driver != nil {
			if driver.Type == r.Type {
				dr, err := driver.New(r.Config)
				if err != nil {
					log.Errorf("Error while creating %s driver:%v", r.Provider, err)
				} else if dr == nil {
					log.Errorf("%s driver is created as nil", r.Provider)
				} else {
					server.addDriver(dr)
				}
			} else {
				log.Errorf("%s named driver is \"%s\" driver in config, but it is actually a \"%s\" driver", r.Provider, r.Type, driver.Type)
			}
		} else {
			log.Errorf("%s named driver not found", r.Provider)
		}
	}
}

func (server *Server) ListenAndServe() {
	port := server.Config.GetPortString()

	//Listen
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server)

	log.Infoln("Listening", port)
	//Serve
	grpcServer.Serve(lis)
}

func NewServer(config *Config) *Server {
	return &Server{
		Config:  config,
		Drivers: make(map[string]drivers.Driver),
	}
}
