package notification

import (
	"errors"
	"notification/drivers"
	pb "notificationpb"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (server *Server) send(ctx context.Context, in *pb.Message, ch chan<- drivers.DriverResult) {
	for _, t := range in.Targets {
		ty := getTargetType(t)
		if driver, ok := server.Drivers[ty]; ok {
			go driver.Send(ctx, in, server.man, ch)
		} else {
			go func() {
				ch <- drivers.DriverResult{Err: errors.New("Driver Not Found"), Type: ty}
			}()
		}
	}
}

func (server *Server) SendMessage(ctx context.Context, in *pb.Message) (*pb.SendMessageResponse, error) {
	if in.Language == "" {
		in.Language = server.Config.DefaultLanguage
	}
	n := len(in.Targets)
	logrus.Debugf("SendMessage with event='%s' and language='%s' to #%d target(s)", in.Event, in.Language, n)
	results := make([]*pb.MessageTargetResponse, 0)
	ch := make(chan drivers.DriverResult, 1)
	go server.send(ctx, in, ch)

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
	if logrus.GetLevel() >= logrus.DebugLevel {
		for _, t := range results {
			logrus.Debugf("SendMessage output[%s]= %s", t.Target, t.Output)
		}
	}
	return pb.NewMessageResponse(results), nil
}

func (server *Server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanResponse, error) {
	logrus.Debugf("Scan")
	ne, err := server.man.Scan()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to scan templates: %v", err)
	}
	return &pb.ScanResponse{Events: ne}, nil
}
