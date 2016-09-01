package notification

import (
	"fmt"
	"net"
	"notification/drivers"
	"notification/template"
	pb "notificationpb"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	Config  *Config
	Drivers map[drivers.NotificationType]drivers.Driver
	man     template.Manager
}

func (server *Server) LoadDrivers() error {
	log.Debugf("server.go: Config is %v", server.Config)
	for _, r := range server.Config.Drivers {
		if driver := drivers.GetDriver(r.Provider); driver != nil {
			if driver.Type == drivers.NotificationType(r.Type) {
				dr, err := driver.New(r.Config)
				if err != nil {
					return fmt.Errorf("failed to create %s driver:%v", r.Provider, err)
				} else if dr == nil {
					return fmt.Errorf("failed to create %s driver, it is created as nil", r.Provider)
				} else {
					server.Drivers[dr.Type()] = dr
				}
			} else {
				return fmt.Errorf("server.go: %s named driver is \"%s\" driver in config, but it is actually a \"%s\" driver", r.Provider, r.Type, driver.Type)
			}
		} else {
			return fmt.Errorf("server.go: %s named driver not found", r.Provider)
		}
	}
	return nil
}

func (server *Server) LoadTemplates() (err error) {
	server.man, err = template.New(server.Config.TemplatePath, server.Config.DefaultLanguage)
	return
}

func (server *Server) ListenAndServe() error {
	port := server.Config.GetPortString()

	//Listen
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server)

	log.Infoln("server.go: Listening", port)
	//Serve
	return grpcServer.Serve(lis)
}

func NewServer(config *Config) *Server {
	return &Server{
		Config:  config,
		Drivers: make(map[drivers.NotificationType]drivers.Driver),
	}
}
