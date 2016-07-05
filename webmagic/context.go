package webmagic

import (
	"net/http"
)

type Context struct {
	Input          *Input
	Output         *Output
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	pathParams     map[string]string
}

func NewContext(rw http.ResponseWriter, r *http.Request) *Context {
	context := new(Context)
	context.Request = r
	context.ResponseWriter = rw
	context.Input = NewInput(r)
	context.Output = NewOutput(rw, r)
	return context
}

func (ctx *Context) FormParam(name string) string {
	if ctx.Request.Form == nil {
		ctx.Request.ParseForm()
	}
	return ctx.Request.Form.Get(name)
}

func (ctx *Context) PathParam(name string) string {
	return ctx.pathParams[name]
}
