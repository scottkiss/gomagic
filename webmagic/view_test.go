package webmagic

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
)

var test_render string = `<!DOCTYPE html>
<html>
  <head>
    <title>test render</title>
  </head>
  <body>
	{{.test}}
  <div>
    <h2>{{.html_text | Html}}</h2>
  </div>
  </body>
</html>
`

func Test_Render(t *testing.T) {
	view := new(View)
	view.Dir = "tmpDir"
	tpl := "index.tpl"
	if err := os.Mkdir("tmpDir", 0777); err != nil {
		t.Fatal(err)
	}
	if f, err := os.Create(filepath.Join("tmpDir", tpl)); err != nil {
		t.Fatal(err)
	} else {
		f.WriteString(test_render)
		f.Close()
	}
	d := make(map[interface{}]interface{})
	d["test"] = "funcmap"
	d["html_text"] = "<a>hello link </a>"
	view.FuncMap = make(template.FuncMap)
	view.FuncMap["Html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	data, err := view.Render(tpl, d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	os.RemoveAll(filepath.Join("tmpDir", tpl))
	os.Remove("tmpDir")
}
