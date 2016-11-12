package soso

func (r *Routes) GET(model string, handler HandlerFunc) {
	r.Handle(model, "get", handler)
}

func (r *Routes) SEARCH(model string, handler HandlerFunc) {
	r.Handle(model, "search", handler)
}

func (r *Routes) CREATE(model string, handler HandlerFunc) {
	r.Handle(model, "create", handler)
}

func (r *Routes) UPDATE(model string, handler HandlerFunc) {
	r.Handle(model, "update", handler)
}

func (r *Routes) DELETE(model string, handler HandlerFunc) {
	r.Handle(model, "delete", handler)
}

func (r *Routes) FLUSH(model string, handler HandlerFunc) {
	r.Handle(model, "flush", handler)
}
