package notification
import pb "notificationpb"

type Job struct {
	RunAt int64
	Data  *pb.Message
}