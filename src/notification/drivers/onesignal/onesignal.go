package onesignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"notification/drivers"
	"notification/template"
	"notificationpb"
	"time"
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
	if len(p.Template) > 0 {
		cts := make(map[string]string)
		if err := json.Unmarshal(p.Template, &cts); err != nil {

		}
	} else {
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
		for _, l := range man.Languages(message.Event, "push.tt") {
			t, err := man.Template(message.Event, l, "push.tt")
			if err != nil {
				return err
			}
			s, err := t.String(data)
			if err != nil {
				return err
			}
			n.Headings[l] = s
		}
	}
	return nil
}

func (o *oneSignalDrv) Send(ctx context.Context, message *notificationpb.Message, man template.Manager, ch chan<- error) {
	p := message.GetPush()
	m := o.notification()
	if uid, ok := message.Tags["user_id"]; ok {
		m.addTagFilter("userid", uid)
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
		ch <- err
		return
	}

	var out bytes.Buffer
	if err := json.NewEncoder(&out).Encode(&m); err != nil {
		ch <- err
		return
	}
	req := httptest.NewRequest("POST", "https://onesignal.com/api/v1/notifications", &out)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", o.authorization))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ch <- errors.New("non ok status code")
		return
	}
	ch <- nil
}
