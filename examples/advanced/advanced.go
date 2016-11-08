package main

import (
	"fmt"

	soso "github.com/happierall/soso-server"
	"github.com/happierall/soso-server/examples/advanced/views"
)

func main() {
	Router := soso.Default()
	Router.HandleRoutes(views.Routes)
	fmt.Println(Router)

	Router.Run(4000)
}
