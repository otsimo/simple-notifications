package sendgrid

import "notification/drivers"


func init() {
	drivers.Register("sendgrid", &drivers.RegisteredDriver{
		Type: drivers.TypeEmail,
		New:newDriver,
	})
}

func newDriver() (drivers.Driver, error) {
	return nil, nil
}
