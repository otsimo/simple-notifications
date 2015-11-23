package notification

import (
	"net"
	"notification/drivers"
	pb "notificationpb"

	"fmt"
	"notification/template"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	Config    *Config
	Drivers   map[string]drivers.Driver
	templates map[string]*template.TemplateGroup
}

func (server *Server) addDriver(dr drivers.Driver) {
	server.Drivers[dr.Type()] = dr
}

func (server *Server) GetTemplateGroup(name string) *template.TemplateGroup {
	return server.templates[name]
}

func (server *Server) LoadDrivers() {
	log.Debugf("server.go: Config is %v", server.Config)

	for _, r := range server.Config.Drivers {
		if driver := drivers.GetDriver(r.Provider); driver != nil {
			if driver.Type == r.Type {
				dr, err := driver.New(r.Config)
				if err != nil {
					log.Errorf("server.go: Error while creating %s driver:%v", r.Provider, err)
				} else if dr == nil {
					log.Errorf("server.go: %s driver is created as nil", r.Provider)
				} else {
					server.addDriver(dr)
				}
			} else {
				log.Errorf("server.go: %s named driver is \"%s\" driver in config, but it is actually a \"%s\" driver", r.Provider, r.Type, driver.Type)
			}
		} else {
			log.Errorf("server.go: %s named driver not found", r.Provider)
		}
	}
}

func (server *Server) LoadTemplates() {
	server.templates = template.SearchTemplates(server.Config.TemplatePath, server.Config.CacheAtStart)
	for _, t := range server.templates {
		log.Debugf("server.go: TemplateGroup[%s] founded", t.Name)
	}
}

func (server *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.SendMessageResponse, error) {
	log.Infof("server.go: SendMessage with event='%s' and language='%s' to #%d target(s)", in.Event, in.Language, len(in.Targets))

	temp := server.GetTemplateGroup(in.Event)
	results := make([]*pb.MessageTargetResponse, 0)

	if temp == nil {
		return pb.NewCustomMessageResponse(pb.ResultEventNotFound,
			fmt.Sprintf("event=%s not found", in.Event), results), nil
	}

	for _, t := range in.Targets {
		ty := getTargetType(t)
		if len(ty) > 0 {
			driver := server.Drivers[ty]
			if driver != nil {
				err := driver.Send(drivers.NewEventData(in.Language, server.Config.DefaultLanguage, t, temp))
				if err != nil {
					results = append(results, &pb.MessageTargetResponse{Type: pb.ResultInternalDriverError, Data: fmt.Sprint(err)})
				} else {
					results = append(results, pb.NewMessageTargetResponse(pb.ResultSuccess, driver.Type(), driver.Name()))
				}
			} else {
				log.Warningf("server.go: %s driver not found", ty)
				results = append(results, pb.NewMessageTargetResponse(pb.ResultDriverNotFound, ty, ""))
			}
		} else {
			log.Warningf("server.go: there is no suitable driver")
			results = append(results, pb.NewMessageTargetResponse(pb.ResultDriverNotFound, "", ""))
		}
	}
	return pb.NewMessageResponse(pb.ResultSuccess, results), nil
}

func (server *Server) ListenAndServe() {
	port := server.Config.GetPortString()

	//Listen
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("server.go: failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server)

	log.Infoln("server.go: Listening", port)
	//Serve
	grpcServer.Serve(lis)
}

func NewServer(config *Config) *Server {
	return &Server{
		Config:  config,
		Drivers: make(map[string]drivers.Driver),
	}
}
