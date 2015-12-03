package webmagic

import (
	"html/template"
)

type View struct {
	//template direction
	Dir string
	//functions map
	FuncMap template.FuncMap
	//cache template
	CacheTpl bool
	//cache template
	templateCache map[string]*template.Template
}

//render template
func (v *View) Render(tpl string, data map[interface{}]interface{}) ([]byte, error) {
	r := new(Render)
	r.root = v.Dir
	r.TplName = tpl
	r.Data = data
	r.FuncMap = v.FuncMap
	r.CacheTpl = v.CacheTpl
	r.templateCache = v.templateCache
	resp, err := r.Build()
	return resp, err
}
