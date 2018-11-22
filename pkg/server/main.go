package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type AwConnection struct {
	conn *websocket.Conn
}

var connections = make(map[string]*AwConnection)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost:8080"
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	server := AwConnection{conn}
	server.serve()
}

type AwSignallingPacket struct {
	SignallingType string      `json:"signallingType"` // ("register" | "relay" | "peerlist")
	SignallingData interface{} `json:"signallingData"` // PeerId | PeerList | AwRelayData;
}

type AwRelaySignallingData struct {
	peerId string
	data   string
}

func (connection *AwConnection) serve() {
	for {
		messageType, receivedJSONmessage, err := connection.conn.ReadMessage()
		if err != nil {
			log.Println(err)
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
			log.Println("registering")
			connection.register(signallingPacket.SignallingData.(string))
			peerList := getPeerIds()
			outgoingSignallingPacket := AwSignallingPacket{
				SignallingType: "peerlist",
				SignallingData: peerList,
			}
			log.Println("outgoingSignallingPacket", outgoingSignallingPacket)
			message, _ := json.Marshal(outgoingSignallingPacket)
			connection.sendMessage(message)
		case "relay":
			log.Println("relaying")
			SignallingData := signallingPacket.SignallingData.(AwRelaySignallingData)
			targetPeerId := SignallingData.peerId
			targetPeerConnection := findConnectionByPeerId(targetPeerId)
			sourcePeerId := findPeerIdForConnection(connection)
			relayData := SignallingData.data
			log.Println("targetPeerId", targetPeerId)
			log.Println("targetPeerConnection", targetPeerConnection)
			log.Println("sourcePeerId", sourcePeerId)
			log.Println("relayData", relayData)
			outgoingSignallingPacket := AwSignallingPacket{
				SignallingType: "relay",
				SignallingData: AwRelaySignallingData{
					peerId: sourcePeerId,
					data:   relayData,
				},
			}
			log.Println("outgoingSignallingPacket", outgoingSignallingPacket)
			message, _ := json.Marshal(outgoingSignallingPacket)
			targetPeerConnection.sendMessage(message)
		}
	}
}

func getPeerIds() []string {
	peerIds := []string{}
	for peerId, _ := range connections {
		peerIds = append(peerIds, peerId)
	}
	return peerIds
}

func findPeerIdForConnection(connection *AwConnection) string {
	var matchingPeerId string
	for peerId, peerConnection := range connections {
		if peerConnection == connection {
			matchingPeerId = peerId
		}
	}
	return matchingPeerId
}

func findConnectionByPeerId(peerId string) *AwConnection {
	return connections[peerId]
}

func (connection *AwConnection) register(peerId string) {
	log.Println("registering peerId ", peerId)
	connections[peerId] = connection
}

func (connection *AwConnection) sendMessage(message []byte) {
	log.Println("sending message", message)
	if err := connection.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
		return
	}
}
