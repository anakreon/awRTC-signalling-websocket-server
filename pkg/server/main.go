package main

import (
	"github.com/anakreon/awrtc-signalling-websocket-server/internal/awconnections"
	"github.com/anakreon/awrtc-signalling-websocket-server/internal/server"
)

func main() {
	awConnections := make(awconnections.AwConnections)
	server := server.WebsocketServer{
		AwConnections: &awConnections,
		Port:          ":8080",
	}
	server.Serve()
}
