package template

import (
	"io"
	"notificationpb"
)

type Template interface {
	Info() (event string, lang string)
	Render(data interface{}, w io.Writer) error
	String(data interface{}) (string, error)
}

type Manager interface {
	Exist(event, lang, suffix string) error
	Languages(event, suffix string) []string
	Template(event, lang, suffix string) (Template, error)
	Scan() ([]*notificationpb.Event, error)
}
