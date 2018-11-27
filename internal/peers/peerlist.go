package peers

type Peerlist map[string]*Peer

func (peerlist *Peerlist) register(peerID string, peer *Peer) {
	(*peerlist)[peerID] = peer
}

func (peerlist *Peerlist) deregister(peerID string) {
	delete(*peerlist, peerID)
}

func (peerlist Peerlist) broadcast(message []byte) {
	for _, peer := range peerlist {
		peer.sendMessage(message)
	}
}

func (peerlist *Peerlist) findPeerByID(peerID string) *Peer {
	return (*peerlist)[peerID]
}

func (peerlist *Peerlist) findIDForPeer(peer *Peer) string {
	var matchingPeerID string
	for peerID, loopPeer := range *peerlist {
		if loopPeer == peer {
			matchingPeerID = peerID
		}
	}
	return matchingPeerID
}

func (peerlist Peerlist) getPeerIDs() []string {
	peerIDs := []string{}
	for peerID := range peerlist {
		peerIDs = append(peerIDs, peerID)
	}
	return peerIDs
}
