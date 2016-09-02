package template

import (
	"errors"
	htmlTemplate "html/template"
	"io/ioutil"
	"notificationpb"
	"path/filepath"
	"strings"
	textTemplate "text/template"
)

var (
	NotFoundError      error = errors.New("not found")
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
	if lang == "" {
		lang = m.defaultLang
	}
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

func (m *manager) Scan() ([]*notificationpb.Event, error) {
	res := make([]*notificationpb.Event, 0)
	fs, err := ioutil.ReadDir(m.dir)
	if err != nil {
		return res, err
	}
	for _, f := range fs {
		if !f.IsDir() {
			continue
		}
		edir := filepath.Join(m.dir, f.Name())
		ef, err := ioutil.ReadDir(edir)
		if err != nil {
			continue
		}
		event := &notificationpb.Event{
			Name:      f.Name(),
			Templates: make([]*notificationpb.Event_Template, 0),
		}
		ts := map[string][]string{}
		for _, ef := range ef {
			if ef.IsDir() {
				continue
			}
			n := ef.Name()
			names := strings.Split(n, ".")
			lindex := 1
			sindex := 2
			if len(names) == 2 {
				lindex = 0
				sindex = 1
			} else if len(names) < 2 {
				continue
			}
			lang := names[lindex]
			suf := names[sindex]
			if _, ok := ts[suf]; !ok {
				ts[suf] = []string{}
			}
			ts[suf] = append(ts[suf], lang)
		}
		for s, l := range ts {
			event.Templates = append(event.Templates, &notificationpb.Event_Template{
				Suffix:    s,
				Languages: l,
			})
		}
		res = append(res, event)
	}
	return res, nil
}

func New(templatesDir, defaultLang string) (Manager, error) {
	return &manager{dir: templatesDir, defaultLang: defaultLang}, nil
}
