package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajayr2512/ptp/peers"
)

func processPeerRequest(conn net.Conn, fileList []string, dataPath string) {
	log.Println("Doing peer server work")

	// TODO: MAJOR: Close connections conn !!
	// TODO: From above maybe no need for delimiters conn closing would be good ?? Check for persistent vs non persistent connections ...
	// TODO: Mem leak checks
	// TODO: Feature: HeartBeats
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

	s := strings.Split(scanner.Text(), ":")
	msgType := s[0]
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

	case "GETFILE":
		fileName := s[1]
		filePath := filepath.Join(dataPath, fileName)
		file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			log.Println(err)
		}

		// TODO: write bytes directly to  conn
		var fileBytes []byte
		if fileBytes, err = ioutil.ReadAll(file); err != nil {
			log.Println(err)
		}
		if _, err := conn.Write(fileBytes); err != nil {
			log.Println(err)
		}
		conn.Close()
		log.Println("Replied with ", string(fileBytes))
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
			go processPeerRequest(con, data, peers.FileLocation(p))
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
	fileHostMap := make(map[string]string)
	fmt.Println("HELP TEXT CLIENT")
loop:
	for {
		fmt.Println("Enter Command: ")
		fmt.Scanf("%s", &in)
		if in == "LEAVE" {
			fmt.Printf("See ya")
			if err := p.Leave(); err != nil {
				log.Fatal(err)
			}
			//stop node server
			close(done)
			break loop
		}

		if strings.HasPrefix(in, "GET") {
			s := strings.Split(in, ":")
			fmt.Println("Getting ", s[1])
			if err := p.GetFile(s[1], fileHostMap); err != nil {
				log.Println(err)
			}
		}

	}

	log.Printf("Peer  registered, data: \n", peerData)
}
