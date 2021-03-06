package peers

import (
	"errors"
	"strconv"
)

func RegisterRequest(peer Peer) (msg []byte, err error) {
	regMsg := []byte("REGISTER\n")
	if peer.hostName == "" || peer.port == 0 {
		return nil, errors.New("peer info corrupt")
	}
	regMsg = append(regMsg, []byte("HOSTNAME:"+peer.hostName+"\n")...)
	regMsg = append(regMsg, []byte("PORT:"+strconv.Itoa(peer.port)+"\n")...)
	regMsg = append(regMsg, []byte("COOKIE:"+strconv.Itoa(peer.cookie)+"\n")...)
	regMsg = append(regMsg, byte('\r'))
	return regMsg, nil
}

func LeaveRequest(peer Peer) (msg []byte, err error) {
	// TODO: Dont need the host port only send cookie
	leaveMsg := []byte("LEAVE\n")
	if peer.cookie == -1 {
		return nil, errors.New("peer not registered")
	}
	leaveMsg = append(leaveMsg, []byte("HOSTNAME:"+peer.hostName+"\n")...)
	leaveMsg = append(leaveMsg, []byte("PORT:"+strconv.Itoa(peer.port)+"\n")...)
	leaveMsg = append(leaveMsg, []byte("COOKIE:"+strconv.Itoa(peer.cookie)+"\n")...)
	leaveMsg = append(leaveMsg, byte('\r'))
	return leaveMsg, nil
}

func ActiveNodesRequest() (msg []byte, err error) {
	nodeMsg := []byte("GETNODES\n\r")
	return nodeMsg, nil
}

func FileListRequest() (msg []byte, err error) {
	listMsg := []byte("GETFILELIST\n\r")
	return listMsg, nil
}

func FileRequest(fileName string) (msg []byte, err error) {
	listMsg := []byte("GETFILE:" + fileName + "\n\r")
	return listMsg, nil
}
