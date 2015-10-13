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
