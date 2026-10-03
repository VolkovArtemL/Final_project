package server

import (
	"log"
	"net/http"
)

const webDir = "./web"

// Run запускает веб-сервер на указанном порту
func Run(port string) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
