package notificationpb

import (
	"fmt"
	"encoding/json"
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

func NewMessageTargetResponse(resultType int32, target, driver string) *MessageTargetResponse {
	return &MessageTargetResponse{
		Type:   resultType,
		Data:   errorMessages[resultType],
		Target: target,
		Driver: driver,
	}
}

func NewMessageResponse(resultType int32, results []*MessageTargetResponse) *SendMessageResponse {
	return &SendMessageResponse{
		Type:    resultType,
		Data:    errorMessages[resultType],
		Results: results,
	}
}

func NewCustomMessageResponse(resultType int32, resultText string, results []*MessageTargetResponse) *SendMessageResponse {
	return &SendMessageResponse{
		Type:    resultType,
		Data:    resultText,
		Results: results,
	}
}

func Map2Str(data map[string]interface{}) string {
	if out, err := json.Marshal(data); err == nil {
		return string(out)
	}
	return "{}"
}