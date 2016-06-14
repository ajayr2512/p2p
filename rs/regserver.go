package main

import (
	"bufio"
	"bytes"
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
			log.Println(i, "From cookie generator")
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
	b := bytes.NewBuffer(msg)
	scanner := bufio.NewScanner(b)
	scanner.Scan()
	msgType := scanner.Text()
	var p peerInfo
	switch msgType {
	case "REGISTER":
		cookie := <-idChan
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
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}

	}
}
