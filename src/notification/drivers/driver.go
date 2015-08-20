package drivers
import "fmt"

type Driver interface{

}

type RegisteredDriver struct {
	New func() (Driver, error)
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