package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ajayr2512/ptp/peers"
)

func processPeerRequest(conn net.Conn, fileList []string) {
	log.Println("Doing peer server work")

	// GET FILE LIST
	b := bufio.NewReader(conn)
	var req []byte
	var err error
	if req, err = b.ReadBytes('\r'); err != nil {
		log.Fatal(err)
	}

	log.Println(string(req))

	// TODO: Directly read req into this
	br := bytes.NewBuffer(req)
	scanner := bufio.NewScanner(br)
	scanner.Scan()
	msgType := scanner.Text()
	log.Println(msgType)

	var reply []byte
	switch msgType {
	case "GETFILELIST":
		for _, i := range fileList {
			reply = append(reply, []byte(i)...)
			reply = append(reply, []byte("\n")...)
		}
		reply = append(reply, byte('\r'))
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}
		log.Println("Replied with: ", string(reply))
	}
	// GET FILE
}

func startNodeServer(p peers.Peer, data []string, done chan struct{}) {
	conChan := make(chan net.Conn)
	ln, err := peers.NewPeerServer(p)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started Node Server")
	var conn net.Conn
	go func() {
		for {
			if conn, err = ln.Accept(); err != nil {
				log.Fatal(err)
			}
			log.Println("Accepted connection")

			conChan <- conn
		}
	}()
loop:
	for {
		select {
		case <-done:
			log.Println("Exiting Peer server")
			break loop
		case con := <-conChan:
			go processPeerRequest(con, data)
		}
	}

}

func main() {
	// Get the cli arguments and populate peer details
	hostName := flag.String("hostname", "localhost", "host name of peer")
	port := flag.Int("port", 85000, "Port Num of peer")
	data := flag.String("data", "", "data associated with peer")

	flag.Parse()
	if *data == "" {
		log.Fatal("Specify data directory for node")
	}

	// Do this on actual file send
	//if err != nil {
	//	log.Fatal("File corrupt")
	//}

	// Check reconnect or new connection
	p, peerData := peers.NewPeer(*port, *hostName, *data)
	log.Println("Starting node server")
	done := make(chan struct{})
	go startNodeServer(p, peerData, done)
	log.Println("Peer Registering with RS")
	if err := p.Register(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Peer  registered, data: \n", peerData, p)
	var in string
	fmt.Println("HELP TEXT CLIENT")
loop:
	for {
		fmt.Scanf("%s", &in)
		if in == "LEAVE" {
			if err := p.Leave(); err != nil {
				log.Fatal(err)
			}
			//stop node server
			close(done)
			break loop
		}

		if strings.HasPrefix(in, "GET") {
			s := strings.Split(in, ":")
			if err := p.GetFile(s[1]); err != nil {
				log.Println(err)
			}
		}

	}

	log.Printf("Peer  registered, data: \n", peerData)
}
