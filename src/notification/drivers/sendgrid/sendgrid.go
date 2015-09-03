package sendgrid

import "notification/drivers"

const SendGridDriverName = "sendgrid"

func init() {
	drivers.Register(SendGridDriverName, &drivers.RegisteredDriver{
		Type: drivers.TypeEmail,
		New: newDriver,
	})
}

func newDriver(config map[string]interface{}) (drivers.Driver, error) {
	d := SendGridDriver{}



	return d, nil
}


type SendGridDriver struct {
}

func (d SendGridDriver)Name() string {
	return SendGridDriverName
}

func (d SendGridDriver)Type() string {
	return drivers.TypeEmail
}