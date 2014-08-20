# About gomagic
gomagic is a middleware magicbox,it is not a framework,but a collection of useful middleware


## Http Magic Useage
```bash
$ go get github.com/scottkiss/gomagic/httpmagic
```

```go
package main

import (
  "github.com/scottkiss/gomagic/httpmagic"
  "log"
  "net/http"
)

func main() {
  ctx := httpmagic.NewContext()
  //handler get request
  //eg. http://localhost:8888/hello/100
  ctx.Get("/hello/:id", handler)
  //handler get request
  ctx.Get("/world/:id", handlerXml)
  //handler post request
  ctx.Post("/post", handlerPost)
  http.Handle("/", ctx)
  http.ListenAndServe(":8888", nil)
}

type User struct {
  Id   string
  Name string
}

//response json
func handler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  id := params.Get(":id")
  log.Println(id)
  user := &User{Id: id, Name: "hello"}
  out := httpmagic.NewOutput(w, r)
  out.Json(user, true)

}

//response xml
func handlerXml(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  id := params.Get(":id")
  log.Println(id)
  user := &User{Id: id, Name: "world"}
  out := httpmagic.NewOutput(w, r)
  out.Xml(user)

}


func handlerPost(w http.ResponseWriter, r *http.Request) {
  in := httpmagic.NewInput(r)
  user := &User{}
  in.ReadJson(user)
  out := httpmagic.NewOutput(w, r)
  out.Json(user, true)

}
```

## License
View the [LICENSE](https://github.com/scottkiss/gomagic/blob/master/LICENSE) file
