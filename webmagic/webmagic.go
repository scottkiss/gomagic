package webmagic

var (
	webmagic *Application
)

func init() {
	webmagic = NewApplication()
}

func Get(pattern string, handler WebHandlerFunc) {
	webmagic.Get(pattern, handler)
}

func Post(pattern string, handler WebHandlerFunc) {
	webmagic.Post(pattern, handler)
}

func Del(pattern string, handler WebHandlerFunc) {
	webmagic.Del(pattern, handler)
}

func Head(pattern string, handler WebHandlerFunc) {
	webmagic.Head(pattern, handler)
}

func Put(pattern string, handler WebHandlerFunc) {
	webmagic.Put(pattern, handler)
}

func Static(pattern string, dir string) {
	webmagic.Static(pattern, dir)
}

func RenderView(tpl string, data map[interface{}]interface{}) ([]byte, error) {
	return webmagic.View().Render(tpl, data)
}

func CacheTpl(cache bool) {
	webmagic.View().CacheTpl = cache
}

func SetViewPath(dir string) {
	webmagic.View().Dir = dir
}

func Run(addr string) {
	webmagic.Run(addr)
}
