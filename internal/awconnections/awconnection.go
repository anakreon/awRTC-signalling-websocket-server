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

type AwSignallingPacket struct {
	SignallingType string `json:"signallingType"` // ("register" | "relay" | "peerlist" | "close")
	SignallingData string `json:"signallingData"` // PeerId | PeerList | AwRelayData;
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
		messageType, receivedJSONmessage, err := connection.conn.ReadMessage()
		if err != nil {
			log.Println("err", err)
			peerID := connection.AwConnections.findPeerIDForConnection(connection)
			connection.AwConnections.deregister(peerID)
			outgoingSignallingPacket := AwSignallingPacket{
				SignallingType: "close",
				SignallingData: peerID,
			}
			log.Println("outgoingSignallingPacket", outgoingSignallingPacket)
			message, _ := json.Marshal(outgoingSignallingPacket)
			connection.AwConnections.broadcast(message)
			return
		}
		if messageType != websocket.TextMessage {
			log.Println("not a text message")
			return
		}
		log.Println("received message", receivedJSONmessage)
		log.Println("received message string", string(receivedJSONmessage))
		signallingPacket := AwSignallingPacket{}
		err2 := json.Unmarshal(receivedJSONmessage, &signallingPacket)
		if err2 != nil {
			log.Println(err2)
			return
		}

		log.Println("signalling packet", signallingPacket)

		switch signallingPacket.SignallingType {
		case "register":
			log.Println("registering user", string(signallingPacket.SignallingData))

			var peerID string
			err2 := json.Unmarshal([]byte(signallingPacket.SignallingData), &peerID)
			if err2 != nil {
				log.Println(err2)
				return
			}
			log.Println("peerid", peerID)
			connection.AwConnections.register(peerID, connection)
			peerList := connection.AwConnections.getPeerIDs()
			peerListJSON, _ := json.Marshal(peerList)
			outgoingSignallingPacket := AwSignallingPacket{
				SignallingType: "peerlist",
				SignallingData: string(peerListJSON),
			}
			log.Println("outgoingSignallingPacket", outgoingSignallingPacket)
			message, _ := json.Marshal(outgoingSignallingPacket)
			connection.sendMessage(message)
		case "relay":
			log.Println("relaying")

			signallingData := AwRelaySignallingData{}
			err2 := json.Unmarshal([]byte(signallingPacket.SignallingData), &signallingData)
			if err2 != nil {
				log.Println(err2)
				return
			}

			targetPeerID := signallingData.PeerId
			targetPeerConnection := connection.AwConnections.findConnectionByPeerID(targetPeerID)
			sourcePeerID := connection.AwConnections.findPeerIDForConnection(connection)
			relayData := signallingData.Data
			log.Println("targetPeerId", targetPeerID)
			log.Println("targetPeerConnection", targetPeerConnection)
			log.Println("sourcePeerId", sourcePeerID)
			log.Println("relayData", relayData)
			outgoingSingnallingData := AwRelaySignallingData{
				PeerId: sourcePeerID,
				Data:   relayData,
			}
			outgoingSingnallingDataJSON, _ := json.Marshal(outgoingSingnallingData)
			outgoingSignallingPacket := AwSignallingPacket{
				SignallingType: "relay",
				SignallingData: string(outgoingSingnallingDataJSON),
			}
			log.Println("outgoingSignallingPacket", outgoingSignallingPacket)
			message, _ := json.Marshal(outgoingSignallingPacket)
			targetPeerConnection.sendMessage(message)
		}
	}
}

func (connection *AwConnection) sendMessage(message []byte) {
	log.Println("sending message", message)
	if err := connection.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
		return
	}
}
