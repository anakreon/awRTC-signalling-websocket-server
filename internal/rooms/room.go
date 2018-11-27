package rooms

import "github.com/anakreon/awrtc-signalling-websocket-server/internal/peers"

type Room struct {
	Peerlist peers.Peerlist
}

func NewRoom() *Room {
	return &Room{
		Peerlist: make(peers.Peerlist),
	}
}
