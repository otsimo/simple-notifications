package notification

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	pb "notificationpb"
)

func (server *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.SendMessageResponse, error) {
	log.Infof("server.go: SendMessage with event='%s' and language='%s' to #%d target(s)", in.Event, in.Language, len(in.Targets))

	results := make([]*pb.MessageTargetResponse, 0)

	ch := make(chan error, 1)
	for _, t := range in.Targets {
		ty := getTargetType(t)
		if driver, ok := server.Drivers[ty]; ok {
			driver.Send(ctx, in, server.man, ch)
		} else {
			log.Warningf("server.go: %s driver not found", ty)
			results = append(results, pb.NewMessageTargetResponse(pb.ResultDriverNotFound, ""))
		}
	}
	return pb.NewMessageResponse(pb.ResultSuccess, results), nil
}
