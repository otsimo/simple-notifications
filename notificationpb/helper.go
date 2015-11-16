package notificationpb

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
