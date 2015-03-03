package httpmagic

import (
	"bytes"
	"github.com/scottkiss/gomagic/utilmagic"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Render struct {
	root    string
	TplName string
	Data    map[interface{}]interface{}
}

func (self *Render) Build() ([]byte, error) {
	outbytes := bytes.NewBufferString("")
	var t *template.Template
	var err error
	t, err = getTemplate(self.root, self.TplName, "")
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

func getTemplate(root, file string, others ...string) (t *template.Template, err error) {
	var filepathAbs string
	filepathAbs = filepath.Join(root, file)
	if exist := utilmagic.FileExists(filepathAbs); !exist {
		panic("can not find template file:" + file)
	}
	data, err := ioutil.ReadFile(filepathAbs)
	if err != nil {
		return nil, err
	}
	t, err = template.New(file).Parse(string(data))
	if err != nil {
		return nil, err
	}
	return t, nil
}
