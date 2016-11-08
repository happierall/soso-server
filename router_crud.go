package soso

func (r *Router) GET(model string, handler HandlerFunc) {
	r.Handle(model, "get", handler)
}

func (r *Router) SEARCH(model string, handler HandlerFunc) {
	r.Handle(model, "search", handler)
}

func (r *Router) CREATE(model string, handler HandlerFunc) {
	r.Handle(model, "create", handler)
}

func (r *Router) UPDATE(model string, handler HandlerFunc) {
	r.Handle(model, "update", handler)
}

func (r *Router) DELETE(model string, handler HandlerFunc) {
	r.Handle(model, "delete", handler)
}

func (r *Router) FLUSH(model string, handler HandlerFunc) {
	r.Handle(model, "flush", handler)
}
