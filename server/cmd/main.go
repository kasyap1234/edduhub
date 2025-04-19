package main

import (
	"eduhub/server/api/app"
	"fmt"
)

func main() {

	setup := app.New().Start()
	if setup != nil {
		fmt.Println("initiated server")
	}
}
