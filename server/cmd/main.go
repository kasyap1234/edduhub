package main

import (
	"eduhub/server/api/app"
	"fmt"
)

func main() {

	setup := app.New()
	err := setup.Start()
	if err != nil {
		fmt.Println("initiated server")
	}
}
