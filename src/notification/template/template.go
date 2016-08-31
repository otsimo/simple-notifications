package template

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	htmltemp "html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	texttemp "text/template"
)

type TemplateType string

const (
	TemplateHtml TemplateType = "html"
	TemplateText TemplateType = "txt"
	TemplateEmailSubject TemplateType = "sub"
	TemplateSms TemplateType = "sms"
	TemplatePush TemplateType = "push"

	TemplateLanguageNone string = ""
)

type tp struct {
	Type     TemplateType
	Language string
	Path     string
	FullName string
	text     *texttemp.Template
	html     *htmltemp.Template
}

type TemplateGroup struct {
	Name      string
	Path      string
	Templates []*tp
}

func (t *TemplateGroup) Find(templateType TemplateType, language string) *tp {
	for _, temp := range t.Templates {
		if temp.Type == templateType && temp.Language == language {
			return temp
		}
	}
	return nil
}

func (t *TemplateGroup) GetText(templateType TemplateType, language, defaultLanguage string, data interface{}) string {
	f := t.Find(templateType, language)
	if f == nil {
		f = t.Find(templateType, defaultLanguage)
		if f == nil {
			log.Errorf("template.go: Template did not found for given language[%s] and type[%v]", language, templateType)
			return ""
		}
	}

	if txt, err := f.PrintText(data); err == nil {
		return txt.String()
	} else {
		log.Errorf("template.go: Error while creating text:%v", err)
	}
	return ""
}

func getTemplateType(fileName string) TemplateType {
	ext := filepath.Ext(fileName)
	if len(ext) < 2 {
		return ""
	}
	return TemplateType(ext[1:])
}

func getLanguage(fileName string) string {
	ext := filepath.Ext(fileName)
	fileName = strings.TrimRight(fileName, ext)
	lang := filepath.Ext(fileName)
	if len(lang) < 2 {
		return TemplateLanguageNone
	}
	return lang[1:]
}

func (t *TemplateGroup) Load() {
	fs, e := ioutil.ReadDir(t.Path)
	if e != nil {
		log.Errorf("template.go: Read sub-template directory '%s' error: %v", t.Path, e)
		return
	}
	log.Debugf("template.go: Loading '%s' template which has %d file(s)", t.Name, len(fs))
	for _, f := range fs {
		if !f.IsDir() {
			name := f.Name()
			templateType := getTemplateType(name)
			language := getLanguage(name)
			log.Debugf("template.go: Template[%s],type='%s',lang='%s' adding to %s", name, templateType, language, t.Name)
			temp := &tp{
				Type:     templateType,
				Language: language,
				FullName: name,
				Path:     filepath.Join(t.Path, name),
			}

			t.Templates = append(t.Templates, temp)
		}
	}
}

func SearchTemplates(root string) map[string]*TemplateGroup {
	res := make(map[string]*TemplateGroup)
	fs, e := ioutil.ReadDir(root)

	if e != nil {
		log.Errorf("template.go: Read template directory '%s' error: %v", root, e)
		return res
	}
	for _, f := range fs {
		if f.IsDir() {
			g := &TemplateGroup{
				Name:      f.Name(),
				Path:      filepath.Join(root, f.Name()),
				Templates: make([]*tp, 0),
			}
			res[g.Name] = g
			g.Load()
		}
	}

	return res
}

func (t *tp) ReadText() ([]byte, error) {
	return ioutil.ReadFile(t.Path)
}

func (t *tp) PrintText(data interface{}) (*bytes.Buffer, error) {
	var doc bytes.Buffer
	var err error

	if t.Type == TemplateHtml {
		err = t.html.Execute(&doc, data)
	} else {
		err = t.text.Execute(&doc, data)
	}

	if err != nil {
		return nil, err
	}
	return &doc, nil //todo potential error
}
