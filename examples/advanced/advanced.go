package main

import (
	soso "../.."
	"./views"
	"github.com/happierall/l"
)

func main() {
	soso.EnableDebug()

	Router := soso.Default()
	Router.HandleRoutes(views.Routes)

	soso.Sessions.OnOpen(func(s soso.Session) {
		l.Log("Session open")
		soso.Sessions.Push(s, "1")
	})

	soso.Sessions.OnClose(func (s soso.Session) {
		l.Log("Session close", s.IsClosed())
	})

	Router.Run(4000)
}
