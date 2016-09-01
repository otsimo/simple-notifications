package notification

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/jsonpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"notification/drivers"
	pb "notificationpb"
	"pipelinepb"
	"strings"
)

const (
	ActionPayloadIsMessage = "payload-is-message"
	ActionConfigIsMessage  = "config-is-message"
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
	logrus.Debugf("SendMessage with event='%s' and language='%s' to #%d target(s)", in.Event, in.Language, len(in.Targets))
	results := make([]*pb.MessageTargetResponse, 0)
	ch := make(chan drivers.DriverResult, 1)
	n := len(in.Targets)

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

func (server *Server) FanOut(*pipelinepb.FlowIn, pipelinepb.Pod_FanOutServer) error {
	return grpc.Errorf(codes.Unimplemented, "FanOut is Unimplemented")
}

func (server *Server) FanIn(pipelinepb.Pod_FanInServer) error {
	return grpc.Errorf(codes.Unimplemented, "FanIn is Unimplemented")
}

func (server *Server) Test(context.Context, *pipelinepb.TestRequest) (*pipelinepb.TestResponse, error) {
	return &pipelinepb.TestResponse{}, nil
}

func (server *Server) Single(ctx context.Context, in *pipelinepb.FlowIn) (*pipelinepb.FlowOut, error) {
	mes := &pb.Message{}
	switch in.Action {
	case ActionPayloadIsMessage:
		if err := jsonpb.UnmarshalString(string(in.Payload), mes); err != nil {
			return nil, err
		}
	case ActionConfigIsMessage:
		if err := jsonpb.UnmarshalString(string(in.Config), mes); err != nil {
			return nil, err
		}
		mes.DataJson = in.Payload
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "unknown action %s", in.Action)
	}
	if mes.Tags == nil {
		mes.Tags = make(map[string]string)
	}
	for k, v := range in.Ids {
		mes.Tags[k] = v
	}
	ch := make(chan drivers.DriverResult, 1)
	n := len(mes.Targets)
	go server.send(ctx, mes, ch)
	errs := []string{}
	for i := 0; i < n; i++ {
		r := <-ch
		if r.Err != nil {
			errs = append(errs, fmt.Sprintf("%v: %v", r.Type, r.Err))
			logrus.Debugf("%v: failed err=%v", r.Type, r.Err)
		} else {
			logrus.Debugf("%v: successful", r.Type)
		}
	}
	if len(errs) > 0 {
		return nil, grpc.Errorf(codes.Internal, "some targets are failed: [ %s ]", strings.Join(errs, " | "))
	}
	return in.Out(), nil
}
