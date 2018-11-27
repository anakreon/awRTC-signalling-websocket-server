package peers

import (
	"encoding/json"
)

type SignallingPacket struct {
	SignallingType string `json:"signallingType"`
	SignallingData string `json:"signallingData"`
}

func NewSignallingPacketPeerList(peerList []string) *SignallingPacket {
	return buildSignallingPacket("peerlist", peerList)
}

func NewSignallingPacketRelay(relaySignallingData RelaySignallingData) *SignallingPacket {
	return buildSignallingPacket("relay", relaySignallingData)
}

func NewSignallingPacketClose(peerID string) *SignallingPacket {
	return buildSignallingPacket("close", peerID)
}

func buildSignallingPacket(signallingType string, signallingData interface{}) *SignallingPacket {
	signallingDataJSON, _ := json.Marshal(signallingData)
	return &SignallingPacket{
		SignallingType: signallingType,
		SignallingData: string(signallingDataJSON),
	}
}

func (packet *SignallingPacket) FromJson(jsonValue []byte) error {
	return json.Unmarshal(jsonValue, packet)
}

func (packet *SignallingPacket) ToJson() (message []byte, err error) {
	return json.Marshal(packet)
}
