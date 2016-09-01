package notificationpb

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func NewMessageTargetResponse(resultType int32, target string) *MessageTargetResponse {
	return &MessageTargetResponse{
		Output: errorMessages[resultType],
		Target: target,
	}
}

func NewMessageResponse(results []*MessageTargetResponse) *SendMessageResponse {
	resp := &SendMessageResponse{
		Results: results,
	}
	failed := []string{}
	for _, r := range results {
		if r.Output != "Success" {
			failed = append(failed, r.Target)
		}
	}
	if len(failed) > 0 {
		resp.Output = fmt.Sprintf("Following targets are failed: %s", strings.Join(failed, " ,"))
	} else {
		resp.Output = "Success"
	}
	return resp
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
