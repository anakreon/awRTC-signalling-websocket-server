package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ConnectionHandler interface {
	Handle(connection *websocket.Conn)
}

type WebsocketServer struct {
	upgrader websocket.Upgrader
	Handler  ConnectionHandler
	Port     string
}

func (server *WebsocketServer) Serve() {
	server.initializeWebsocketUpgrader()
	http.HandleFunc("/", server.websocketHandler)
	log.Fatal(http.ListenAndServe(server.Port, nil))
}

func (server *WebsocketServer) initializeWebsocketUpgrader() {
	server.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Host == "localhost:8080"
		},
	}
}

func (server *WebsocketServer) websocketHandler(writer http.ResponseWriter, request *http.Request) {
	connection, err := server.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	server.Handler.Handle(connection)
}
