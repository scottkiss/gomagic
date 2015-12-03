package webmagic

import (
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	HTTPMETHOD = map[string]string{
		"DELETE": "DELETE",
		"HEAD":   "HEAD",
		"GET":    "GET",
		"POST":   "POST",
		"PUT":    "PUT",
	}
)

type route struct {
	method  string
	regexp  *regexp.Regexp
	params  map[int]string
	handler http.HandlerFunc
}

type Application struct {
	routes  []*route
	filters []http.HandlerFunc
	view    *View
}

func NewApplication() *Application {
	app := new(Application)
	app.view = new(View)
	app.view.templateCache = make(map[string]*template.Template)
	return app
}

func (app *Application) View() *View {
	return app.view
}

func (app *Application) Get(pattern string, handler http.HandlerFunc) {
	app.RegisterRoute(HTTPMETHOD["GET"], pattern, handler)
}

func (app *Application) Del(pattern string, handler http.HandlerFunc) {
	app.RegisterRoute(HTTPMETHOD["DELETE"], pattern, handler)
}

func (app *Application) Post(pattern string, handler http.HandlerFunc) {
	app.RegisterRoute(HTTPMETHOD["POST"], pattern, handler)
}

func (app *Application) Head(pattern string, handler http.HandlerFunc) {
	app.RegisterRoute(HTTPMETHOD["HEAD"], pattern, handler)
}

func (app *Application) Put(pattern string, handler http.HandlerFunc) {
	app.RegisterRoute(HTTPMETHOD["PUT"], pattern, handler)
}

func (app *Application) RegisterRoute(method string, pattern string, handler http.HandlerFunc) {
	subpath := strings.Split(pattern, "/")
	params := make(map[int]string)
	j := 0
	for i, p := range subpath {
		if strings.HasPrefix(p, ":") {
			expr := "([^/]+)"
			if index := strings.Index(p, "("); index != -1 {
				expr = p[index:]
				p = p[:index]
			}
			params[j] = p
			subpath[i] = expr
			j++
		}

	}

	pattern = strings.Join(subpath, "/")
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
		return
	}
	route := &route{}
	route.method = method
	route.params = params
	route.handler = handler
	route.regexp = regex

	app.routes = append(app.routes, route)
}

func (app *Application) Filter(filter http.HandlerFunc) {
	app.filters = append(app.filters, filter)
}

func (app *Application) FilterParam(param string, filter http.HandlerFunc) {
	if !strings.HasPrefix(param, ":") {
		param = ":" + param
	}
	app.Filter(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get(param)
		if len(p) > 0 {
			filter(w, r)
		}
	})
}

func (app *Application) Static(pattern string, dir string) {
	pattern = pattern + "(.+)"
	app.RegisterRoute(HTTPMETHOD["GET"], pattern, func(rw http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)
		path = filepath.Join(dir, path)
		http.ServeFile(rw, r, path)
	})
}

func (app *Application) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	w := &responseWrite{writer: rw}
	for _, route := range app.routes {
		if r.Method != route.method {
			continue
		}
		if !route.regexp.MatchString(requestPath) {
			continue
		}

		matches := route.regexp.FindStringSubmatch(requestPath)
		if len(matches[0]) != len(requestPath) {
			continue
		}
		if len(route.params) > 0 {
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.params[i], match)
			}
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery

		}

		for _, filter := range app.filters {
			filter(w, r)
			if w.started {
				return
			}
		}

		route.handler(w, r)
		break
	}

	if w.started == false {
		http.NotFound(w, r)
	}
}

type responseWrite struct {
	writer  http.ResponseWriter
	status  int
	started bool
}

func (rw *responseWrite) Header() http.Header {
	return rw.writer.Header()
}

func (rw *responseWrite) Write(b []byte) (int, error) {
	rw.started = true
	return rw.writer.Write(b)
}

func (rw *responseWrite) WriteHeader(code int) {
	rw.status = code
	rw.started = true
	rw.writer.WriteHeader(code)
}

func (app *Application) Run(addr string) {
	if addr == "" {
		panic("input address invalid")
	}
	println("http server run at " + addr)
	e := http.ListenAndServe(addr, app)
	panic(e)
}
