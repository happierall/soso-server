package main

import (
	soso "github.com/happierall/soso-server"
	"github.com/happierall/soso-server/examples/advanced/views"
)

func main() {
	Router := soso.Default()
	Router.HandleRoutes(views.Routes)

	Router.Run(4000)
}
