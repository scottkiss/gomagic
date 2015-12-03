package webmagic

import (
	"bytes"
	"github.com/scottkiss/gomagic/utilmagic"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Render struct {
	root          string
	TplName       string
	Data          map[interface{}]interface{}
	FuncMap       template.FuncMap
	CacheTpl      bool
	templateCache map[string]*template.Template
}

func (self *Render) Build() ([]byte, error) {
	outbytes := bytes.NewBufferString("")
	var t *template.Template
	var err error
	t, err = self.getTemplate()
	if err != nil {
		log.Panic("getTemplate err:", err)
		return nil, err
	}
	err = t.ExecuteTemplate(outbytes, self.TplName, self.Data)
	if err != nil {
		log.Panic("template Execute error:", err)
		return nil, err
	}

	content, _ := ioutil.ReadAll(outbytes)
	return content, nil

}

func (self *Render) getTemplate() (t *template.Template, err error) {
	if self.CacheTpl {
		if self.templateCache[self.TplName] != nil {
			return self.templateCache[self.TplName], nil
		}
	}
	var filepathAbs string
	filepathAbs = filepath.Join(self.root, self.TplName)
	if exist := utilmagic.FileExists(filepathAbs); !exist {
		panic("can not find template file:" + self.TplName)
	}
	data, err := ioutil.ReadFile(filepathAbs)
	if err != nil {
		return nil, err
	}
	t = template.New(self.TplName)
	t.Funcs(self.FuncMap)
	t.Parse(string(data))
	if err != nil {
		return nil, err
	}
	if self.CacheTpl {
		self.templateCache[self.TplName] = t
	}
	return t, nil
}
