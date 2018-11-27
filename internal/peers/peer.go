package peers

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Peer struct {
	Peerlist *Peerlist
	conn     *websocket.Conn
}

type RelaySignallingData struct {
	PeerId string `json:"peerId"`
	Data   string `json:"data"`
}

func (peer *Peer) Handle(conn *websocket.Conn) {
	peer.conn = conn
	peer.serve()
}

func (peer *Peer) serve() {
	for {
		messageType, message, error := peer.conn.ReadMessage()
		if error != nil {
			peer.close()
			return
		}
		if messageType != websocket.TextMessage {
			log.Println("not a text message")
			return
		}

		incomingSignallingPacket := SignallingPacket{}
		incomingSignallingPacket.FromJson(message)

		switch incomingSignallingPacket.SignallingType {
		case "register":
			peer.register(incomingSignallingPacket.SignallingData)
		case "relay":
			peer.relay(incomingSignallingPacket.SignallingData)
		}
	}
}

func (peer *Peer) register(signallingData string) {
	peerID := getPeerID(signallingData)
	peer.Peerlist.register(peerID, peer)
	peer.sendPeerList()
}

func getPeerID(signallingData string) (peerID string) {
	json.Unmarshal([]byte(signallingData), &peerID)
	return
}

func (peer *Peer) sendPeerList() {
	peerList := peer.Peerlist.getPeerIDs()
	signallingPacket := NewSignallingPacketPeerList(peerList)
	peer.sendSignallingPacket(signallingPacket)
}

func (peer *Peer) sendSignallingPacket(signallingPacket *SignallingPacket) {
	message, _ := signallingPacket.ToJson()
	peer.sendMessage(message)
}

func (peer *Peer) broadcastSignallingPacket(signallingPacket *SignallingPacket) {
	message, _ := signallingPacket.ToJson()
	peer.Peerlist.broadcast(message)
}

func (peer *Peer) sendSignallingPacketToPeer(peerID string, signallingPacket *SignallingPacket) {
	peerConnection := peer.Peerlist.findPeerByID(peerID)
	peerConnection.sendSignallingPacket(signallingPacket)
}

func (peer *Peer) relay(signallingData string) {
	incomingRelaySignallingData := getRelaySignallingData(signallingData)
	outgoingRelaySignallingData := peer.buildRelaySignallingData(incomingRelaySignallingData.Data)
	outgoingSignallingPacket := NewSignallingPacketRelay(outgoingRelaySignallingData)
	peer.sendSignallingPacketToPeer(incomingRelaySignallingData.PeerId, outgoingSignallingPacket)
}

func (peer *Peer) buildRelaySignallingData(data string) RelaySignallingData {
	currentPeerID := peer.Peerlist.findIDForPeer(peer)
	return RelaySignallingData{
		PeerId: currentPeerID,
		Data:   data,
	}
}

func getRelaySignallingData(signallingData string) RelaySignallingData {
	relaySignallingData := RelaySignallingData{}
	json.Unmarshal([]byte(signallingData), &relaySignallingData)
	return relaySignallingData
}

func (peer *Peer) close() {
	peerID := peer.Peerlist.findIDForPeer(peer)
	peer.Peerlist.deregister(peerID)
	signallingPacket := NewSignallingPacketClose(peerID)
	peer.broadcastSignallingPacket(signallingPacket)
}

func (peer *Peer) sendMessage(message []byte) {
	if err := peer.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
		return
	}
}
