package template_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"notification/template"

	log "github.com/Sirupsen/logrus"
)

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareTestDirectoryWithLanguage(root string, templates []string) {
	os.Mkdir(root, os.ModePerm)
	for _, t := range templates {
		a := filepath.Join(root, t)
		data := "Hello World!! " + t

		err := os.Mkdir(a, os.ModePerm)
		check(err)
		err = ioutil.WriteFile(filepath.Join(a, t+".en.sms"), []byte(data), os.ModePerm)
		check(err)
		err = ioutil.WriteFile(filepath.Join(a, t+".tr.sms"), []byte(data), os.ModePerm)
		check(err)
	}
	//	err := filepath.Walk(filepath.Join(root, TestDirectory), visit)
	//	fmt.Printf("filepath.Walk() returned %v\n", err)
}

func prepareTestDirectoryWithoutLanguage(root string, templates []string) {
	os.Mkdir(root, os.ModePerm)

	for _, t := range templates {
		a := filepath.Join(root, t)
		data := "Hello World!! " + t

		err := os.Mkdir(a, os.ModePerm)
		check(err)
		err = ioutil.WriteFile(filepath.Join(a, t+".sms"), []byte(data), os.ModePerm)
		check(err)
		err = ioutil.WriteFile(filepath.Join(a, t+".push"), []byte(data), os.ModePerm)
		check(err)
	}
}

func TestSearchTemplates(t *testing.T) {
	root := filepath.Join(filepath.Dir(os.Args[0]), ".test1")
	templates := []string{"abc", "xyz"}

	prepareTestDirectoryWithLanguage(root, templates)

	l := template.SearchTemplates(root, false)

	for _, n := range templates {
		temp := l[n]
		if temp == nil {
			t.Fatalf("templateGroup[%s] not found", n)
		}
		if len(temp.Templates) != 2 {
			t.Fatalf("templateGroup[%s] don't have 2 templates", n)
		}
		trsms := temp.Find(template.TemplateSms, "tr")
		if trsms == nil {
			t.Fatal("TR Sms template is nil")
		}
		ensms := temp.Find(template.TemplateSms, "en")
		if ensms == nil {
			t.Fatal("EN Sms template is nil")
		}
	}
	os.Remove(root)
}

func TestSearchTemplatesWithoutLanguage(t *testing.T) {
	root := filepath.Join(filepath.Dir(os.Args[0]), ".test2")
	templates := []string{"asd", "fgh"}
	prepareTestDirectoryWithoutLanguage(root, templates)

	l := template.SearchTemplates(root, false)

	for _, n := range templates {
		temp := l[n]
		if temp == nil {
			t.Fatalf("templateGroup[%s] not found", n)
		}
		if len(temp.Templates) != 2 {
			t.Fatalf("templateGroup[%s] don't have 2 templates", n)
		}

		sms := temp.Find(template.TemplateSms, template.TemplateLanguageNone)
		if sms == nil {
			t.Fatal("Sms template is nil")
		}
		push := temp.Find(template.TemplatePush, template.TemplateLanguageNone)
		if push == nil {
			t.Fatal("Push template is nil")
		}
	}
	os.Remove(root)
}

func init() {
	log.SetLevel(log.DebugLevel)
}