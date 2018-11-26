package awconnections

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type AwConnection struct {
	AwConnections *AwConnections
	conn          *websocket.Conn
}

type AwRelaySignallingData struct {
	PeerId string `json:"peerId"`
	Data   string `json:"data"`
}

func (connection *AwConnection) Handle(conn *websocket.Conn) {
	connection.conn = conn
	connection.serve()
}

func (connection *AwConnection) serve() {
	for {
		messageType, message, error := connection.conn.ReadMessage()
		if error != nil {
			connection.close()
			return
		}
		if messageType != websocket.TextMessage {
			log.Println("not a text message")
			return
		}

		incomingSignallingPacket := AwSignallingPacket{}
		incomingSignallingPacket.FromJson(message)

		switch incomingSignallingPacket.SignallingType {
		case "register":
			connection.register(incomingSignallingPacket.SignallingData)
		case "relay":
			connection.relay(incomingSignallingPacket.SignallingData)
		}
	}
}

func (connection *AwConnection) register(signallingData string) {
	peerID := getPeerID(signallingData)
	connection.AwConnections.register(peerID, connection)
	connection.sendPeerList()
}

func getPeerID(signallingData string) (peerID string) {
	json.Unmarshal([]byte(signallingData), &peerID)
	return
}

func (connection *AwConnection) sendPeerList() {
	peerList := connection.AwConnections.getPeerIDs()
	signallingPacket := NewSignallingPacketPeerList(peerList)
	connection.sendSignallingPacket(signallingPacket)
}

func (connection *AwConnection) sendSignallingPacket(signallingPacket *AwSignallingPacket) {
	log.Println("sending packet", signallingPacket)
	message, _ := signallingPacket.ToJson()
	connection.sendMessage(message)
}

func (connection *AwConnection) broadcastSignallingPacket(signallingPacket *AwSignallingPacket) {
	message, _ := signallingPacket.ToJson()
	connection.AwConnections.broadcast(message)
}

func (connection *AwConnection) sendSignallingPacketToPeer(peerID string, signallingPacket *AwSignallingPacket) {
	peerConnection := connection.AwConnections.findConnectionByPeerID(peerID)
	peerConnection.sendSignallingPacket(signallingPacket)
}

func (connection *AwConnection) relay(signallingData string) {
	incomingRelaySignallingData := getRelaySignallingData(signallingData)
	outgoingRelaySignallingData := connection.buildRelaySignallingData(incomingRelaySignallingData.Data)
	outgoingSignallingPacket := NewSignallingPacketRelay(outgoingRelaySignallingData)
	connection.sendSignallingPacketToPeer(incomingRelaySignallingData.PeerId, outgoingSignallingPacket)
}

func (connection *AwConnection) buildRelaySignallingData(data string) AwRelaySignallingData {
	currentPeerID := connection.AwConnections.findPeerIDForConnection(connection)
	return AwRelaySignallingData{
		PeerId: currentPeerID,
		Data:   data,
	}
}

func getRelaySignallingData(signallingData string) AwRelaySignallingData {
	relaySignallingData := AwRelaySignallingData{}
	json.Unmarshal([]byte(signallingData), &relaySignallingData)
	return relaySignallingData
}

func (connection *AwConnection) close() {
	peerID := connection.AwConnections.findPeerIDForConnection(connection)
	connection.AwConnections.deregister(peerID)
	signallingPacket := NewSignallingPacketClose(peerID)
	connection.broadcastSignallingPacket(signallingPacket)
}

func (connection *AwConnection) sendMessage(message []byte) {
	if err := connection.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
		return
	}
}
