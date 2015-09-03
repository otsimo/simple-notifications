package drivers
import (
	"fmt"
)

type Driver interface {
	Name() string
	Type() string
}

const (
	TypeEmail string = "email"
	TypeSms string = "sms"
	TypePush string = "push"
	TypeScheduler string = "scheduler"
)

type RegisteredDriver struct {
	Type string
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
	drivers[name] = rd;
	return nil
}

func GetDrivers() map[string]string {
	drives := make(map[string]string, 0)

	for name, d := range (drivers) {
		drives[name] = d.Type
	}
	return drives
}

func GetDriver(name string) *RegisteredDriver {
	return drivers[name]
}