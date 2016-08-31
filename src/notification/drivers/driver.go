package drivers

import (
	"fmt"
	"notification/template"
	"notificationpb"
	"golang.org/x/net/context"
)

type Driver interface {
	Name() string
	Type() NotificationType
	Send(ctx context.Context, message *notificationpb.Message, man template.Manager, ch chan<- error)
}
type NotificationType string

const (
	TypeEmail   NotificationType = "email"
	TypeSms     NotificationType = "sms"
	TypePush    NotificationType = "push"
	TypeUnknown NotificationType = ""
)

type RegisteredDriver struct {
	Type NotificationType
	New  func(map[string]interface{}) (Driver, error)
}

var drivers map[string]*RegisteredDriver

func init() {
	drivers = make(map[string]*RegisteredDriver)
}

func Register(name string, rd *RegisteredDriver) error {
	if _, ext := drivers[name]; ext {
		return fmt.Errorf("Name already registered %s", name)
	}
	drivers[name] = rd
	return nil
}

func GetDrivers() map[string]NotificationType {
	drives := make(map[string]NotificationType, 0)

	for name, d := range drivers {
		drives[name] = d.Type
	}
	return drives
}

func GetDriver(name string) *RegisteredDriver {
	return drivers[name]
}
