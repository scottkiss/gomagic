package httpmagic

import (
	"html/template"
)

type View struct {
	//template direction
	Dir string
	//functions map
	FucMap template.FuncMap
}

//render template
func (v *View) Render(tpl string, data map[interface{}]interface{}) ([]byte, error) {
	r := &Render{v.Dir, tpl, data}
	resp, err := r.Build()
	return resp, err
}
