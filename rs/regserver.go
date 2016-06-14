package main

import (
	"bufio"
	"bytes"
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

var peerDict = make(map[int]peerInfo)

//var idChan chan int

// Deciding on whether it is the same host joining
// or a new host is based on hostName/port combo
// TO DO: Add new msg join which enables re-joining
// on peer control.
func main() {
	idChan := make(chan int)

	// Starting cookie generator
	go func() {
		i := 1
		for {
			idChan <- i
			log.Println(i, "From cookie generator")
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
	log.Println("Starting Registration Server")
	conChan := make(chan net.Conn)
	go func() {
		for {
			var conn net.Conn
			if conn, err = ln.Accept(); err != nil {
				log.Fatal(err)
			}
			log.Println("Accepted connection")
			conChan <- conn
		}
	}()

	log.Println("Waiting for connections")
	for {
		conn := <-conChan
		log.Println("Processing connection")
		b := bufio.NewReader(conn)
		var req []byte
		if req, err = b.ReadBytes('\r'); err != nil {
			log.Fatal(err)
		}

		log.Println(string(req))
		go processRequest(req, conn, idChan)
	}

}

func processRequest(msg []byte, conn net.Conn, idChan chan int) {
	log.Println("Populating reply")
	b := bytes.NewBuffer(msg)
	scanner := bufio.NewScanner(b)
	scanner.Scan()
	msgType := scanner.Text()
	log.Println(msgType)
	var p peerInfo
	switch msgType {
	case "REGISTER":
		cookie := <-idChan
		log.Println(cookie)
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ":")
			switch s[0] {
			case "HOSTNAME":
				p.hostName = s[1]
			case "PORT":
				p.port, _ = strconv.Atoi(s[1])
			}
		}

		p.ttl = 7200
		p.flag = true
		peerDict[cookie] = p

		reply := []byte("STATUS:PASS\nCOOKIE:" + strconv.Itoa(cookie))
		reply = append(reply, byte('\r'))
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}

	}
}
