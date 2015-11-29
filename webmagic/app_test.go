package webmagic

import (
	"../webmagic"
	"net/http"
	"testing"
)

func Test_Run(t *testing.T) {
	defer func() {
		e := recover()
		if e != nil {
			t.Fatal(e)
		}
	}()
	app := webmagic.NewApplication()
	app.Get("/hello/:id", handler)
	app.Run(":8888")

}

func handler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	out := webmagic.NewOutput(w, r)
	out.Html([]byte("hello webmagic , id is " + id))
}
