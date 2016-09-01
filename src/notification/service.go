package notification

import (
	"notification/drivers"
	pb "notificationpb"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (server *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.SendMessageResponse, error) {
	log.Debugf("service.go: SendMessage with event='%s' and language='%s' to #%d target(s)", in.Event, in.Language, len(in.Targets))
	results := make([]*pb.MessageTargetResponse, 0)
	ch := make(chan drivers.DriverResult, 1)
	n := 0
	for _, t := range in.Targets {
		ty := getTargetType(t)
		if driver, ok := server.Drivers[ty]; ok {
			n++
			go driver.Send(ctx, in, server.man, ch)
		} else {
			log.Warningf("service.go: %s driver not found", ty)
			results = append(results, pb.NewMessageTargetResponse(pb.ResultDriverNotFound, string(ty)))
		}
	}
	for i := 0; i < n; i++ {
		r := <-ch
		resp := &pb.MessageTargetResponse{
			Target: string(r.Type),
			Output: "Success",
		}
		if r.Err != nil {
			resp.Output = r.Err.Error()
		}
		results = append(results, resp)
	}
	return pb.NewMessageResponse(results), nil
}

func (server *Server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanResponse, error) {
	log.Debugf("Scan")
	ne, err := server.man.Scan()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to scan templates: %v", err)
	}
	return &pb.ScanResponse{Events: ne}, nil
}
