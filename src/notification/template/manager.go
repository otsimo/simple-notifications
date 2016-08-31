package template

import (
	"bytes"
	htmlTemplate "html/template"
	"io"
	textTemplate "text/template"
)

type manager struct {
	dir         string
	defaultLang string
}

func (m *manager) Exist(event, lang, suffix string) error {
	return nil
}
func (m *manager) Languages(event, suffix string) []string {
	return nil
}
func (m *manager) Template(event, lang, suffix string) (Template, error) {
	return nil, nil
}

type Template interface {
	Render(data interface{}, w io.Writer) error
	String(data interface{}) (string, error)
}

type Manager interface {
	Exist(event, lang, suffix string) error
	Languages(event, suffix string) []string
	Template(event, lang, suffix string) (Template, error)
}

func New(templatesDir, defaultLang string) (Manager, error) {
	return &manager{dir: templatesDir, defaultLang: defaultLang}, nil
}

type txtTemplate struct {
	t *textTemplate.Template
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

type hTemplate struct {
	t *htmlTemplate.Template
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
