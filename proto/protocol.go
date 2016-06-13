package proto

import (
	"errors"
	"hrank/ptp/peers"
	"strconv"
)

func RegisterMessage(peer peers.Peer) (msg []byte, err error) {
	regMsg = []byte("REGISTER\n")
	if peer.hostName == nil || peer.port == 0 {
		return nil, errors.New("peer info corrupt")
	}
	regMsg.append(regMsg, []byte("HOSTNAME : "+peer.hostName+"\n"))
	regMsg.append(regMsg, []byte("PORT : "+strconv.Itoa(peer.port)))
	return regMsg, nil
}
