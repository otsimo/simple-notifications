package notification

import (
	"notification/drivers"
	pb "notificationpb"
)

func getTargetType(t *pb.Target) drivers.NotificationType {
	switch t.GetBackend().(type) {
	case *pb.Target_Email:
		return drivers.TypeEmail
	case *pb.Target_Push:
		return drivers.TypePush
	case *pb.Target_Sms:
		return drivers.TypeSms
	default:
		return drivers.TypeUnknown
	}
}
