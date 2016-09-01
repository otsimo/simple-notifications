package template

import (
	"bytes"
	"errors"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	textTemplate "text/template"
)

var (
	NotFoundError error = errors.New("not found")
	EventNotExistError error = errors.New("event not exist")
)

type manager struct {
	dir         string
	defaultLang string
}

func (m *manager) templateFile(event, lang, suffix string) (string, error) {
	edir := filepath.Join(m.dir, event)
	fs, e := ioutil.ReadDir(edir)
	if e != nil {
		return "", EventNotExistError
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		n := f.Name()
		names := strings.Split(n, ".")
		lindex := 1
		sindex := 2
		if len(names) == 2 {
			lindex = 0
			sindex = 1
		} else if len(names) < 2 {
			continue
		}
		if names[lindex] == lang && strings.HasSuffix(names[sindex], suffix) {
			return filepath.Join(edir, f.Name()), nil
		}
	}
	return "", NotFoundError
}

func (m *manager) Exist(event, lang, suffix string) (err error) {
	_, err = m.templateFile(event, lang, suffix)
	return
}

func (m *manager) Languages(event, suffix string) []string {
	edir := filepath.Join(m.dir, event)
	fs, e := ioutil.ReadDir(edir)
	res := []string{}
	if e != nil {
		return res
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		n := f.Name()
		names := strings.Split(n, ".")
		lindex := 1
		sindex := 2
		if len(names) == 2 {
			lindex = 0
			sindex = 1
		} else if len(names) < 2 {
			continue
		}
		if strings.HasSuffix(names[sindex], suffix) {
			res = append(res, names[lindex])
		}
	}
	return res
}

func (m *manager) Template(event, lang, suffix string) (Template, error) {
	t, err := m.templateFile(event, lang, suffix)
	if err != nil {
		lang = m.defaultLang
		t, err = m.templateFile(event, m.defaultLang, suffix)
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(t)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(t, "html") {
		t, err := htmlTemplate.New("").Parse(string(b))
		if err != nil {
			return nil, err
		}
		return &hTemplate{t: t, event: event, lang: lang}, nil
	} else {
		t, err := textTemplate.New("").Parse(string(b))
		if err != nil {
			return nil, err
		}
		return &txtTemplate{t: t, event: event, lang: lang}, nil
	}
}

type Template interface {
	Info() (event string, lang string)
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
