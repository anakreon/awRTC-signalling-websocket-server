package awconnections

import (
	"encoding/json"
)

type AwSignallingPacket struct {
	SignallingType string `json:"signallingType"`
	SignallingData string `json:"signallingData"`
}

func NewSignallingPacketPeerList(peerList []string) *AwSignallingPacket {
	return buildSignallingPacket("peerlist", peerList)
}

func NewSignallingPacketRelay(relaySignallingData AwRelaySignallingData) *AwSignallingPacket {
	return buildSignallingPacket("relay", relaySignallingData)
}

func NewSignallingPacketClose(peerID string) *AwSignallingPacket {
	return buildSignallingPacket("close", peerID)
}

func buildSignallingPacket(signallingType string, signallingData interface{}) *AwSignallingPacket {
	signallingDataJSON, _ := json.Marshal(signallingData)
	return &AwSignallingPacket{
		SignallingType: signallingType,
		SignallingData: string(signallingDataJSON),
	}
}

func (packet *AwSignallingPacket) FromJson(jsonValue []byte) error {
	return json.Unmarshal(jsonValue, packet)
}

func (packet *AwSignallingPacket) ToJson() (message []byte, err error) {
	return json.Marshal(packet)
}
