package RS

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

type peerInfo struct {
	ttl      int
	flag     bool
	hostName string
	port     int
}

var peerDict map[int]peerInfo
var idChan chan int

// Deciding on whether it is the same host joining
// or a new host is based on hostName/port combo
// TO DO: Add new msg join which enables re-joining
// on peer control.
func main() {

	log.Println("Starting Registration Server")
	// Starting cookie generator
	go func() {
		i := 1
		for {
			idChan <- i
			i = i + 1
		}
	}()

	defer func() {
		log.Println("Exiting RS")
	}()
	ln, err := net.Listen("tcp", ":60000")
	if err != nil {
		log.Fatal(err)
	}
	conChan := make(chan net.Conn)
	go func() {
		for {
			var conn net.Conn
			if conn, err = ln.Accept(); err != nil {
				log.Fatal(err)
			}
			conChan <- conn
		}
	}()

	for {
		conn := <-conChan
		var msg []byte

		if msg, err = ioutil.ReadAll(conn); err != nil {
			log.Fatal(err)
		}

		go processRequest(msg, conn)

	}

}

func processRequest(msg []byte, conn net.Conn) {
	scanner := bufio.NewScanner(msg)
	scanner.Scan()
	msgType := scanner.Text()
	switch msgType {
	case "REGISTER":
		cookie := <-idChan
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ":")
			switch s[0] {
			case "HOSTNAME":
				peerInfo.hostName = s[1]
			case "PORT":
				peerInfo.port = strconv.Atoi(s[1])
			}
		}

		peerInfo.ttl = 7200
		peerInfo.flag = true
		peerDict[cookie] = peerInfo

		reply := []byte("STATUS:PASS\nCOOKIE:" + stconv.Itoa(cookie))
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}

	}
}
