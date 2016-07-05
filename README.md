# About gomagic
gomagic is a middleware magicbox,it is not a framework,but a collection of useful middleware


## Web Magic Useage
```bash
$ go get github.com/scottkiss/gomagic/webmagic
```

```go
package main

import (
  "github.com/scottkiss/gomagic/webmagic"
  "log"
)

func main() {
  app := webmagic.Application()
  //handler get request
  //eg. http://localhost:8888/hello/100
  app.Get("/hello/:id", handler)
  //handler get request
  app.Get("/world/:id", handlerXml)
  //handler post request
  app.Post("/post", handlerPost)
  app.Run(":8888")
}

type User struct {
  Id   string
  Name string
}

//response json
func handler(ctx *webmagic.Context) {
  id := ctx.PathParam("id")
  log.Println(id)
  user := &User{Id: id, Name: "hello"}
  ctx.Output.Json(user, true)
}

//response xml
func handlerXml(ctx *webmagic.Context) {
  id := ctx.PathParam("id")
  log.Println(id)
  user := &User{Id: id, Name: "world"}
  ctx.Output.Xml(user)
}


func handlerPost(ctx *webmagic.Context) {
  user := &User{}
  ctx.Input.ReadJson(user)
  ctx.Output.Json(user, true)
}
```

## License
View the [LICENSE](https://github.com/scottkiss/gomagic/blob/master/LICENSE) file
