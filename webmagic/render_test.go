package webmagic

import (
	"os"
	"path/filepath"
	"testing"
)

var index string = `<!DOCTYPE html>
<html>
  <head>
    <title>welcome template</title>
  </head>
  <body>
	{{.say}}
  </body>
</html>
`

func Test_Build(t *testing.T) {
	dir := "tmpDir"
	file := "index.tpl"
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}
	if f, err := os.Create(filepath.Join(dir, file)); err != nil {
		t.Fatal(err)
	} else {
		f.WriteString(index)
		f.Close()
	}
	d := make(map[interface{}]interface{})
	d["say"] = "hello"
	r := &Render{dir, file, d, nil, false, nil}
	data, err := r.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	os.RemoveAll(filepath.Join(dir, file))
	os.RemoveAll(dir)

}
