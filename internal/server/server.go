package server

import (
	"log"
	"net/http"

	"github.com/anakreon/awRTC-signalling-websocket-server/internal/peers"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
	Peerlist *peers.Peerlist
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
	conn, err := server.upgrader.Upgrade(writer, request, nil)
	peer := peers.Peer{
		Peerlist: server.Peerlist,
	}
	if err != nil {
		log.Println(err)
		return
	}

	peer.Handle(conn)
}
