package peers

import (
	"errors"
	"strconv"
)

func RegisterMessage(peer Peer) (msg []byte, err error) {
	regMsg := []byte("REGISTER\n")
	if peer.hostName == "" || peer.port == 0 {
		return nil, errors.New("peer info corrupt")
	}
	regMsg = append(regMsg, []byte("HOSTNAME : "+peer.hostName+"\n")...)
	regMsg = append(regMsg, []byte("PORT : "+strconv.Itoa(peer.port))...)
	return regMsg, nil
}
