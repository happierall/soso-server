package main

import (
	"fmt"

	soso "github.com/happierall/soso-server"
)

func main() {
	Router := soso.Default()

	Router.CREATE("message", func(m *soso.Msg) {
		fmt.Println(m.RequestMap)

		m.Success(map[string]interface{}{
			"id": 1,
		})
	})

	Router.Run(4000)
}
