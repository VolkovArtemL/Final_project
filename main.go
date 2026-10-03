package main

import (
	"log"
	"os"

	"github.com/VolkovArtemL/Final_project/pkg/db"
	"github.com/VolkovArtemL/Final_project/pkg/server"
)

const defaultPort = "7540"
const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := server.Run(port); err != nil {
		log.Fatal(err)
	}
}
