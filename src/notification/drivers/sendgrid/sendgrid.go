package sendgrid

import (
	"notification/drivers"
	pb "notificationpb"

	log "github.com/Sirupsen/logrus"
	sendgrid "github.com/sendgrid/sendgrid-go"
)

const SendGridDriverName = "sendgrid"

func init() {
	drivers.Register(SendGridDriverName, &drivers.RegisteredDriver{
		Type: drivers.TypeEmail,
		New:  newDriver,
	})
}

func newDriver(config map[string]interface{}) (drivers.Driver, error) {
	d := &SendGridDriver{}

	if config["apiUser"] != nil {
		d.Client = sendgrid.NewSendGridClient(config["apiUser"].(string), config["apiKey"].(string))
	} else {
		d.Client = sendgrid.NewSendGridClientWithApiKey(config["apiKey"].(string))
	}
	return d, nil
}

type SendGridDriver struct {
	Client           *sendgrid.SGClient
	DefaultFromEmail string
	DefaultFromName  string
}

func (d SendGridDriver) Name() string {
	return SendGridDriverName
}

func (d SendGridDriver) Type() string {
	return drivers.TypeEmail
}

func (d *SendGridDriver) Send(in *pb.Message, t *pb.Target) error {
	log.Infoln("Sending mail via", SendGridDriverName)

	m := t.GetEmail()

	message := sendgrid.NewMail()

	message.SetFrom(m.FromEmail)
	message.AddTos(m.ToEmail)
	message.SetSubject(m.Subject)

	if len(m.ToName) > 0 {
		message.AddToNames(m.ToName)
	}
	if len(m.FromName) > 0 {
		message.SetFromName(m.FromName)
	}
	//TODO Read And Parse
	message.SetText("Hello World")

	r := d.Client.Send(message)
	if r == nil {
		log.Infoln("Email sent!")
	} else {
		log.Errorln(r)
	}
	return r
}
