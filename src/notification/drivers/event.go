package drivers

import (
	"notification/template"
	pb "notificationpb"
	"encoding/json"
)

type EventData struct {
	Language        string
	DefaultLanguage string
	Target          *pb.Target
	TemplateGroup   *template.TemplateGroup
}

func (e *EventData)GetEmailData() map[string]interface{} {
	data := make(map[string]interface{})
	email := e.Target.GetEmail()
	if email != nil {
		json.Unmarshal([]byte(email.DataJson), &data)
	}
	return data
}

func (e *EventData)GetSmsData() map[string]interface{} {
	data := make(map[string]interface{})
	sms := e.Target.GetSms()
	if sms != nil {
		json.Unmarshal([]byte(sms.DataJson), &data)
	}
	return data
}

func (e *EventData)GetPushData() map[string]interface{} {
	data := make(map[string]interface{})
	push := e.Target.GetPush()
	if push != nil {
		json.Unmarshal([]byte(push.DataJson), &data)
	}
	return data
}

func (e *EventData) GetHtml(data interface{}) string {
	return e.TemplateGroup.GetText(template.TemplateHtml, e.Language, e.DefaultLanguage, data)
}

func (e *EventData) GetText(t template.TemplateType, data interface{}) string {
	return e.TemplateGroup.GetText(t, e.Language, e.DefaultLanguage, data)
}

func NewEventData(lang, defaultLang string, target *pb.Target, temp *template.TemplateGroup) EventData {
	return EventData{
		Language:lang,
		DefaultLanguage:defaultLang,
		Target:target,
		TemplateGroup:temp,
	}
}