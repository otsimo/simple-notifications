package sendgrid

import (
	"encoding/json"
	"errors"
	"fmt"
	"notification/drivers"
	"notification/template"
	"notificationpb"

	"github.com/Sirupsen/logrus"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/net/context"
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

	if ak, ok := config["apiKey"]; ok {
		d.ApiKey, ok = ak.(string)
		if !ok {
			return nil, errors.New("apiKey must be string")
		}
	} else {
		return nil, errors.New("missing apiKey")
	}

	if sender := config["defaultFrom"]; sender != nil {
		d.DefaultFromEmail = sender.(string)
	}
	if sender := config["defaultFromName"]; sender != nil {
		d.DefaultFromName = sender.(string)
	}
	logrus.WithField("driver", SendGridDriverName).Errorf("initialized")
	return d, nil
}

type SendGridDriver struct {
	DefaultFromEmail string
	DefaultFromName  string
	ApiKey           string
}

func (d SendGridDriver) Name() string {
	return SendGridDriverName
}

func (d SendGridDriver) Type() drivers.NotificationType {
	return drivers.TypeEmail
}

func templateString(data interface{}, customTemplate []byte, message *notificationpb.Message, man template.Manager, suffix string, html bool) (str string, set bool, err error) {
	set = true
	var temp template.Template
	if len(customTemplate) > 0 {
		temp, err = template.NewTemplate(string(customTemplate), html)
	} else {
		if man.Exist(message.Event, message.Language, suffix) == nil {
			temp, err = man.Template(message.Event, message.Language, suffix)
		} else {
			return "", false, nil
		}
	}
	if err != nil {
		return "", true, err
	}
	str, err = temp.String(data)
	return
}

func (d *SendGridDriver) Send(ctx context.Context, message *notificationpb.Message, man template.Manager, ch chan<- drivers.DriverResult) {
	m := message.GetEmail()
	email := new(mail.SGMailV3)
	p := mail.NewPersonalization()
	var fromName string
	if len(m.FromName) > 0 {
		fromName = m.FromName
	} else {
		fromName = d.DefaultFromName
	}
	if len(m.FromEmail) > 0 {
		email.SetFrom(mail.NewEmail(fromName, m.FromEmail))
	} else {
		email.SetFrom(mail.NewEmail(fromName, d.DefaultFromEmail))
	}
	addName := len(m.ToEmail) == len(m.ToName)
	for i, e := range m.ToEmail {
		if addName {
			p.AddTos(mail.NewEmail(m.ToName[i], e))
		} else {
			p.AddTos(mail.NewEmail("", e))
		}
	}
	for _, e := range m.Cc {
		p.AddCCs(mail.NewEmail("", e))
	}
	for _, e := range m.Bcc {
		p.AddBCCs(mail.NewEmail("", e))
	}
	if len(m.ReplyTo) > 0 {
		email.SetReplyTo(mail.NewEmail("", m.ReplyTo))
	}
	data := make(map[string]interface{})
	if len(message.DataJson) > 0 {
		if err := json.Unmarshal(message.DataJson, &data); err != nil {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: err}
			return
		}
	}
	for k, v := range message.Tags {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	var err error
	email.Subject, _, err = templateString(data, m.TemplateSub, message, man, "sub", false)
	if err != nil {
		ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: errors.New("failed to create subject text")}
		return
	}
	var html string
	var set bool
	var text string
	text, set, err = templateString(data, m.TemplateTxt, message, man, "txt", false)
	if set {
		if err != nil {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: errors.New("failed to create text content")}
			return
		}
		email.AddContent(mail.NewContent("text/plain", text))
	}
	html, set, err = templateString(data, m.TemplateHtml, message, man, "html", true)
	if set {
		if err != nil {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: errors.New("failed to create html content")}
			return
		}
		email.AddContent(mail.NewContent("text/html", html))
	}
	if message.ScheduleAt > 0 {
		email.SendAt = int(message.ScheduleAt)
	}
	email.AddPersonalizations(p)
	request := sendgrid.GetRequest(d.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(email)
	select {
	case <-ctx.Done():
		ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: ctx.Err()}
	default:
		r, err := sendgrid.API(request)
		if err != nil {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: err}
			return
		}
		if r.StatusCode >= 300 {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: fmt.Errorf(r.Body)}
		} else {
			ch <- drivers.DriverResult{Type: drivers.TypeEmail, Err: nil}
		}
	}
}
