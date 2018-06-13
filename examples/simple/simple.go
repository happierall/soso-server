package main

import (
	"fmt"

	soso "../.."
)

func main() {

	soso.EnableDebug()

	Router := soso.Default()

	Router.CREATE("user", func(m *soso.Msg) {
		m.Success(map[string]interface{}{
			"id": 1,
		})
	})

	fmt.Println(Router.Run(4000))
}
