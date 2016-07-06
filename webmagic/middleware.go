package webmagic

type Middleware struct {
	Name    string
	Handler MiddlewareHandlerFunc
	next    *Middleware
}

func (middle Middleware) Next(ctx *Context) {
	if next := middle.next; next != nil {
		next.Handler(ctx, next)
	}
}

type MiddlewareHandlerFunc func(*Context, *Middleware)
