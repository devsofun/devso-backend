package main

import (
	"devso-backend/routes"
)

func main() {
	router := routes.SetupRouter()
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
