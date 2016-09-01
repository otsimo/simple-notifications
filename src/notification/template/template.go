package template

import (
	"bytes"
	htmlTemplate "html/template"
	"io"
	textTemplate "text/template"
)

type txtTemplate struct {
	t     *textTemplate.Template
	event string
	lang  string
}

func (t *txtTemplate) Render(data interface{}, w io.Writer) error {
	return t.t.Execute(w, data)
}

func (t *txtTemplate) String(data interface{}) (string, error) {
	var b bytes.Buffer
	if err := t.t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
func (t *txtTemplate) Info() (string, string) {
	return t.event, t.lang
}

type hTemplate struct {
	t     *htmlTemplate.Template
	event string
	lang  string
}

func (t *hTemplate) Render(data interface{}, w io.Writer) error {
	return t.t.Execute(w, data)
}

func (t *hTemplate) String(data interface{}) (string, error) {
	var b bytes.Buffer
	if err := t.t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (t *hTemplate) Info() (string, string) {
	return t.event, t.lang
}

func NewTemplate(str string, isHtml bool) (Template, error) {
	if isHtml {
		t, err := htmlTemplate.New("").Parse(str)
		if err != nil {
			return nil, err
		}
		return &hTemplate{t: t}, nil
	} else {
		t, err := textTemplate.New("").Parse(str)
		if err != nil {
			return nil, err
		}
		return &txtTemplate{t: t}, nil
	}
}
