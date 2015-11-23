package sendgrid

import (
	"notification/drivers"
	"notification/template"

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

	if sender := config["defaultFrom"]; sender != nil {
		d.DefaultFromEmail = sender.(string)
	}
	if sender := config["defaultFromName"]; sender != nil {
		d.DefaultFromName = sender.(string)
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

func (d *SendGridDriver) Send(data drivers.EventData) error {
	m := data.Target.GetEmail()
	md := data.GetEmailData()

	message := sendgrid.NewMail()

	if len(m.FromEmail) > 0 {
		message.SetFrom(m.FromEmail)
	} else {
		message.SetFrom(d.DefaultFromEmail)
	}

	message.AddTos(m.ToEmail)
	message.AddCcs(m.Cc)
	message.AddBccs(m.Bcc)
	if len(m.ReplyTo) > 0 {
		message.SetReplyTo(m.ReplyTo)
	}
	if len(m.ToName) > 0 {
		message.AddToNames(m.ToName)
	}
	if len(m.FromName) > 0 {
		message.SetFromName(m.FromName)
	} else {
		message.SetFromName(d.DefaultFromName)
	}

	if txt := data.GetHtml(md); len(txt) > 0 {
		message.SetHTML(txt)
	}

	if txt := data.GetText(template.TemplateText, md); len(txt) > 0 {
		message.SetText(txt)
	}

	if txt := data.GetText(template.TemplateEmailSubject, md); len(txt) > 0 {
		message.SetSubject(txt)
	} else {
		message.SetSubject(m.Subject)
	}

	r := d.Client.Send(message)
	if r != nil {
		log.Errorln(r)
	}
	return r
}