package main

import (
	"github.com/anakreon/awrtc-signalling-websocket-server/internal/awconnections"
	"github.com/anakreon/awrtc-signalling-websocket-server/internal/server"
)

func main() {
	awConnections := make(awconnections.AwConnections)
	awConnection := awconnections.AwConnection{
		AwConnections: &awConnections,
	}
	server := server.WebsocketServer{
		Handler: &awConnection,
		Port:    ":8080",
	}
	server.Serve()
}
