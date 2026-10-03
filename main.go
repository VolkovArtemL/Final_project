package main

import (
	"log"
	"os"

	"github.com/VolkovArtemL/Final_project/pkg/server"
)

const defaultPort = "7540"

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := server.Run(port); err != nil {
		log.Fatal(err)
	}
}
