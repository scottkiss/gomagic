package webmagic

import (
	"net/http"
)

var (
	webmagic *Application
)

func init() {
	webmagic = NewApplication()
}

func Get(pattern string, handler http.HandlerFunc) {
	webmagic.Get(pattern, handler)
}

func Post(pattern string, handler http.HandlerFunc) {
	webmagic.Post(pattern, handler)
}

func Del(pattern string, handler http.HandlerFunc) {
	webmagic.Del(pattern, handler)
}

func Head(pattern string, handler http.HandlerFunc) {
	webmagic.Head(pattern, handler)
}

func Put(pattern string, handler http.HandlerFunc) {
	webmagic.Put(pattern, handler)
}

func Static(pattern string, dir string) {
	webmagic.Static(pattern, dir)
}

func Run(addr string) {
	webmagic.Run(addr)
}
