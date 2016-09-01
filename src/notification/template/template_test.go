package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareTestDirectory(root string, templates []testTemplate) {
	data := "Hello World!! "
	for _, t := range templates {
		a := filepath.Join(root, t.event)
		err := os.MkdirAll(a, os.ModePerm)
		check(err)
		for _, f := range t.files {
			err = ioutil.WriteFile(filepath.Join(a, f), []byte(data), os.ModePerm)
			check(err)
		}
	}
}

type testTemplate struct {
	event string
	files []string
}

func TestSearchTemplates(t *testing.T) {
	tests := []struct {
		defaultLang string
		templates   []testTemplate
		tevent      string
		tlang       string
		tsuffix     string
		revent      string
		rlang       string
		err         error
	}{
		{
			defaultLang: "en",
			templates: []testTemplate{{
				event: "abc",
				files: []string{"a.en.txt", "a.tr.txt", "a.en.push", "a.tr.push"},
			}, {
				event: "xyz",
				files: []string{"a.en.txt", "a.tr.txt", "a.en.push", "a.tr.push"},
			}},
			tevent:  "abc",
			tlang:   "tr",
			tsuffix: "push",
			revent:  "abc",
			rlang:   "tr",
			err:     nil,
		},
		{
			defaultLang: "en",
			templates: []testTemplate{{
				event: "abc",
				files: []string{"a.en.txt", "a.en.push"},
			}, {
				event: "xyz",
				files: []string{"a.en.txt", "a.en.push"},
			}},
			tevent:  "xyz",
			tlang:   "tr",
			tsuffix: "push",
			revent:  "xyz",
			rlang:   "en",
			err:     nil,
		},
		{
			defaultLang: "en",
			templates: []testTemplate{{
				event: "abc",
				files: []string{"a.en.txt", "a.en.push"},
			}},
			tevent:  "xyz",
			tlang:   "tr",
			tsuffix: "push",
			revent:  "xyz",
			rlang:   "en",
			err:     EventNotExistError,
		}, {
			defaultLang: "en",
			templates: []testTemplate{{
				event: "abc",
				files: []string{"en.txt", "tr.txt", "en.push", "tr.push"},
			}, {
				event: "xyz",
				files: []string{"en.txt", "tr.txt", "en.push", "tr.push"},
			}},
			tevent:  "abc",
			tlang:   "tr",
			tsuffix: "push",
			revent:  "abc",
			rlang:   "tr",
			err:     nil,
		}, {
			defaultLang: "en",
			templates: []testTemplate{{
				event: "abc",
				files: []string{"en.xhtml", "tr.xhtml", "en.push", "tr.push"},
			}, {
				event: "xyz",
				files: []string{"en.html", "tr.html", "en.push", "tr.push"},
			}},
			tevent:  "abc",
			tlang:   "tr",
			tsuffix: "html",
			revent:  "abc",
			rlang:   "tr",
			err:     nil,
		},
	}

	for i, test := range tests {
		root := filepath.Join(os.TempDir(), "not", "search", fmt.Sprintf("test%d", i))
		prepareTestDirectory(root, test.templates)
		man, err := New(root, test.defaultLang)
		if err != nil {
			t.Fatalf("case %d: %v", i, err)
		}
		temp, err := man.Template(test.tevent, test.tlang, test.tsuffix)
		if test.err != nil && err == nil {
			t.Fatalf("case %d: wants err=%v got=nil", i, test.err)
		} else if test.err == nil && err != nil {
			t.Fatalf("case %d: wants err=nil got=%v", i, err)
		} else if test.err != nil && err != nil {
			if test.err.Error() != err.Error() {
				t.Fatalf("case %d: wants err=%v got=%v", i, test.err, err)
			}
		}
		if test.err != nil {
			os.Remove(root)
			continue
		}
		ev, lang := temp.Info()
		if ev != test.revent {
			t.Fatalf("case %d: event wants=%v got=%v", i, test.revent, ev)
		}
		if lang != test.rlang {
			t.Fatalf("case %d: lang wants=%v got=%v", i, test.rlang, lang)
		}
		os.Remove(root)
	}
}

func init() {
	log.SetLevel(log.DebugLevel)
}
