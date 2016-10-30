package soso

func (r *Router) CREATE(data_type string, handler HandlerFunc) {
	r.Handle(data_type, "create", handler)
}

func (r *Router) RETRIEVE(data_type string, handler HandlerFunc) {
	r.Handle(data_type, "retrieve", handler)
}

func (r *Router) UPDATE(data_type string, handler HandlerFunc) {
	r.Handle(data_type, "update", handler)
}

func (r *Router) DELETE(data_type string, handler HandlerFunc) {
	r.Handle(data_type, "delete", handler)
}

func (r *Router) FLUSH(data_type string, handler HandlerFunc) {
	r.Handle(data_type, "flush", handler)
}
