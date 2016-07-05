package webmagic

import (
	"testing"
)

func Test_Run(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			t.Fatal(e)
		}
	}()
	app := NewApplication()
	app.Get("/hello/:id", handler)
	app.Run(":8888")

}

func handler(ctx *Context) {
	id := ctx.PathParam("id")
	ctx.Output.Html([]byte("hello webmagic , id is " + id))
}
