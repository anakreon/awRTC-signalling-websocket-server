package main

import (
	"github.com/anakreon/awRTC-signalling-websocket-server/internal/peers"
	"github.com/anakreon/awRTC-signalling-websocket-server/internal/server"
)

func main() {
	peerlist := make(peers.Peerlist)
	server := server.WebsocketServer{
		Peerlist: &peerlist,
		Port:     ":8080",
	}
	server.Serve()
}
