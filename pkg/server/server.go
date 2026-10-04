package server

import (
	"log"
	"net/http"

	"github.com/VolkovArtemL/Final_project/pkg/api"
)

const webDir = "./web"

func Run(port string) error {
	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
