package notificationpb

const (
	ResultSuccess             int32 = 0
	ResultDriverNotFound      int32 = 1001
	ResultInternalDriverError int32 = 1002
	ResultTemplateNotFound    int32 = 1003
	ResultEventNotFound       int32 = 1004
)

var errorMessages map[int32]string

func init() {
	errorMessages = map[int32]string{
		ResultSuccess:             "Success",
		ResultDriverNotFound:      "Driver Not Found",
		ResultInternalDriverError: "Driver Send Error",
		ResultTemplateNotFound:    "Template Not Found",
		ResultEventNotFound:       "Event Not Found",
	}
}
