package sendgrid

import "notification/drivers"


func init() {
	drivers.Register("sendgrid", &drivers.RegisteredDriver{
		New:newDriver,
	})
}

func newDriver() (drivers.Driver, error) {
	return nil, nil
}
