package notificationpb

import (
	"encoding/json"
	"fmt"
)

func NewEmailTarget(email *Email) *Target {
	return &Target{
		Backend: &Target_Email{
			Email: email,
		},
	}
}

func NewSmsTarget(sms *Sms) *Target {
	return &Target{
		Backend: &Target_Sms{
			Sms: sms,
		},
	}
}

func NewPushTarget(push *Push) *Target {
	return &Target{
		Backend: &Target_Push{
			Push: push,
		},
	}
}

func NewTargets(targets ...interface{}) []*Target {
	r := make([]*Target, 0)
	for _, opt := range targets {
		switch v := opt.(type) {
		case *Email:
			r = append(r, &Target{
				Backend: &Target_Email{
					Email: opt.(*Email),
				},
			})
		case *Sms:
			r = append(r, &Target{
				Backend: &Target_Sms{
					Sms: opt.(*Sms),
				},
			})
		case *Push:
			r = append(r, &Target{
				Backend: &Target_Push{
					Push: opt.(*Push),
				},
			})
		default:
			fmt.Printf("unknown notification target %v", v)
		}
	}
	return r
}

func NewMessageTargetResponse(resultType int32, driver string) *MessageTargetResponse {
	return &MessageTargetResponse{
		Output: errorMessages[resultType],
		Driver: driver,
	}
}

func NewMessageResponse(resultType int32, results []*MessageTargetResponse) *SendMessageResponse {
	return &SendMessageResponse{
		Output:  errorMessages[resultType],
		Results: results,
	}
}

func NewCustomMessageResponse(output string, results []*MessageTargetResponse) *SendMessageResponse {
	return &SendMessageResponse{
		Output:  output,
		Results: results,
	}
}

func Map2Str(data map[string]interface{}) []byte {
	if out, err := json.Marshal(data); err == nil {
		return out
	}
	return []byte("{}")
}

func (m *Message) GetEmail() *Email {
	for _, t := range m.Targets {
		if e := t.GetEmail(); e != nil {
			return e
		}
	}
	return nil
}

func (m *Message) GetSms() *Sms {
	for _, t := range m.Targets {
		if e := t.GetSms(); e != nil {
			return e
		}
	}
	return nil
}

func (m *Message) GetPush() *Push {
	for _, t := range m.Targets {
		if e := t.GetPush(); e != nil {
			return e
		}
	}
	return nil
}
