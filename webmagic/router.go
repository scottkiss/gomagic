package webmagic

import (
	"regexp"
)

type Route struct {
	method  string
	regexp  *regexp.Regexp
	params  map[int]string
	handler WebHandlerFunc
}

type WebHandlerFunc func(ctx *Context)
