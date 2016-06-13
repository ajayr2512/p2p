package RS

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type peerInfo struct {
	ttl  int
	flag bool
}

// Deciding on whether it is the same host joining
// or a new host is based on hostName/port combo
// TO DO: Add new msg join which enables re-joining
// on peer control.
func main() {
	log.Println("Starting Registration Server")
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
	// FINISH UP
	scanner := bufio.NewScanner(msg)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		switch s[0] {
		case "REGISTER":
			// Handle this
		}
	}
}
