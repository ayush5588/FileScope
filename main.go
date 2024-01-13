package main

import (
	"log"

	"github.com/ayush5588/FileScope/internal/router"
)

const (
	portNumber = ":8080"
)

func main() {
	r := router.SetupRouter()
	err := r.Run(portNumber)
	if err != nil {
		log.Fatal("error in starting the router...", err)
	}
}
