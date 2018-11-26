package awconnections

type AwConnections map[string]*AwConnection

func (connections *AwConnections) register(peerID string, connection *AwConnection) {
	(*connections)[peerID] = connection
}

func (connections *AwConnections) deregister(peerID string) {
	delete(*connections, peerID)
}

func (connections AwConnections) broadcast(message []byte) {
	for _, connection := range connections {
		connection.sendMessage(message)
	}
}

func (connections *AwConnections) findConnectionByPeerID(peerID string) *AwConnection {
	return (*connections)[peerID]
}

func (connections *AwConnections) findPeerIDForConnection(connection *AwConnection) string {
	var matchingPeerID string
	for peerID, peerConnection := range *connections {
		if peerConnection == connection {
			matchingPeerID = peerID
		}
	}
	return matchingPeerID
}

func (connections AwConnections) getPeerIDs() []string {
	peerIDs := []string{}
	for peerID := range connections {
		peerIDs = append(peerIDs, peerID)
	}
	return peerIDs
}
