package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/loukikbhandari/mqr"
)

func main() {
	connection, err := mqr.OpenConnection("handler", "tcp", "localhost:6379", 2, nil)
	if err != nil {
		panic(err)
	}

	http.Handle("/overview", NewHandler(connection))
	fmt.Printf("Handler listening on http://localhost:3333/overview\n")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		panic(err)
	}
}

type Handler struct {
	connection mqr.Connection
}

func NewHandler(connection mqr.Connection) *Handler {
	return &Handler{connection: connection}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	layout := request.FormValue("layout")
	refresh := request.FormValue("refresh")

	queues, err := handler.connection.GetOpenQueues()
	if err != nil {
		panic(err)
	}

	stats, err := handler.connection.CollectStats(queues)
	if err != nil {
		panic(err)
	}

	log.Printf("queue stats\n%s", stats)
	fmt.Fprint(writer, stats.GetHtml(layout, refresh))
}
