package onesignal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notification/drivers"
	"notification/template"
	"notificationpb"
	"time"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

const DriverName = "onesignal"

func init() {
	drivers.Register(DriverName, &drivers.RegisteredDriver{
		Type: drivers.TypePush,
		New:  newDriver,
	})
}

func newDriver(config map[string]interface{}) (drivers.Driver, error) {
	d := &oneSignalDrv{}
	if ai, ok := config["appID"]; ok {
		d.appID = ai.(string)
	}
	if a, ok := config["authorization"]; ok {
		d.authorization = a.(string)
	}
	logrus.WithField("driver", DriverName).Info("initialized")
	return d, nil
}

type oneSignalDrv struct {
	authorization string
	appID         string
}

func (o *oneSignalDrv) Name() string {
	return DriverName
}
func (o *oneSignalDrv) Type() drivers.NotificationType {
	return drivers.TypePush
}

type filter struct {
	Field    string `json:"field"`
	Key      string `json:"key,omitempty"`
	Relation string `json:"relation,omitempty"`
	Value    string `json:"value,omitempty"`
}

func newTagFilter(key, value, relation string) filter {
	return filter{
		Field:    "tag",
		Key:      key,
		Value:    value,
		Relation: relation,
	}
}

type notification struct {
	AppID            string            `json:"app_id"`
	Contents         map[string]string `json:"contents"`
	Headings         map[string]string `json:"headings,omitempty"`
	Data             map[string]string `json:"data,omitempty"`
	SendAfter        string            `json:"send_after,omitempty"`
	Filters          []filter          `json:"filters,omitempty"`
	IncludePlayerIDs []string          `json:"include_player_ids,omitempty"`
}

func (n *notification) addTagFilter(key, value string) {
	n.Filters = append(n.Filters, newTagFilter(key, value, "="))
}

func (o *oneSignalDrv) notification() notification {
	return notification{
		AppID:    o.appID,
		Filters:  make([]filter, 0),
		Contents: make(map[string]string),
		Headings: make(map[string]string),
	}
}

type output struct {
	ID         string   `json:"id,omitempty"`
	Recipients int      `json:"recipients,omitempty"`
	Error      []string `json:"errors,omitempty"`
}

func putContent(n *notification, message *notificationpb.Message, man template.Manager) error {
	p := message.GetPush()
	langs := man.Languages(message.Event, "push")
	data := make(map[string]interface{})
	if err := json.Unmarshal(message.DataJson, &data); err != nil {
		return err
	}
	for k, v := range message.Tags {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	for _, l := range langs {
		t, err := man.Template(message.Event, l, "push")
		if err != nil {
			return err
		}
		s, err := t.String(data)
		if err != nil {
			return err
		}
		n.Contents[l] = s
	}
	for _, l := range man.Languages(message.Event, "tit") {
		t, err := man.Template(message.Event, l, "tit")
		if err != nil {
			return err
		}
		s, err := t.String(data)
		if err != nil {
			return err
		}
		n.Headings[l] = s
	}
	if len(p.Template) > 0 {
		cts := make(map[string]string)
		if err := json.Unmarshal(p.Template, &cts); err != nil {
			if len(langs) == 0 {
				t, e := template.NewTemplate(string(p.Template), false)
				if e != nil {
					return e
				}
				s, _ := t.String(data)
				n.Contents["en"] = s
				n.Contents[message.Language] = s
				return nil
			}
		} else {
			for k, v := range cts {
				t, e := template.NewTemplate(string(v), false)
				if e != nil {
					continue
				}
				n.Contents[k], _ = t.String(data)
			}
			if _, ok := n.Contents["en"]; !ok {
				n.Contents["en"] = ""
			}
			return nil
		}
	}
	return nil
}

func (o *oneSignalDrv) Send(ctx context.Context, message *notificationpb.Message, man template.Manager, ch chan<- drivers.DriverResult) {
	p := message.GetPush()
	m := o.notification()
	if uid, ok := message.Tags["user_id"]; ok {
		m.addTagFilter("userid", uid)
	} else if len(p.To) > 0 {
		m.addTagFilter("userid", p.To[0])
	}

	if env, ok := message.Tags["env"]; ok {
		m.addTagFilter("env", env)
	}
	if len(p.To) > 0 {
		m.IncludePlayerIDs = p.To
	}
	if message.ScheduleAt > 0 {
		m.SendAfter = time.Unix(message.ScheduleAt, 0).UTC().String()
	}
	if err := putContent(&m, message, man); err != nil {
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: err}
		return
	}

	var out bytes.Buffer
	if err := json.NewEncoder(&out).Encode(&m); err != nil {
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: fmt.Errorf("encode message json error, %v", err)}
		return
	}
	req, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", &out)
	if err != nil {
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: err}
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", o.authorization))
	req.Header.Set("Content-Type", "application/json")
	select {
	case <-ctx.Done():
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: ctx.Err()}
		return
	default:
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: err}
		return
	}
	if resp.StatusCode != http.StatusOK {
		ch <- drivers.DriverResult{Type: drivers.TypePush, Err: fmt.Errorf("status code: %d", resp.StatusCode)}
		return
	}
	ch <- drivers.DriverResult{Type: drivers.TypePush, Err: nil}
}
