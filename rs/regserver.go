package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
)

// TODO: Need flag ??
type peerInfo struct {
	ttl      int
	flag     bool
	hostName string
	port     int
}

var activeDict = make(map[int]peerInfo)

// Moved cookie logic + Connect/Re-connect to clients

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

	// TODO: move the reading from connection to process request
	// and dont pass req ..
	log.Println("Waiting for connections")
	for {
		conn := <-conChan
		log.Println("Processing connection")
		b := bufio.NewReader(conn)
		var req []byte
		if req, err = b.ReadBytes('\r'); err != nil {
			log.Println(err)
			continue
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
	var cookie int
	switch msgType {
	case "REGISTER":
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ":")
			switch s[0] {
			case "HOSTNAME":
				p.hostName = s[1]
			case "PORT":
				p.port, _ = strconv.Atoi(s[1])
			case "COOKIE":
				cookie, _ = strconv.Atoi(s[1])
			}
		}
		var reply []byte
		if cookie == -1 {
			cookie = <-idChan
			reply = []byte("STATUS:NEW\nCOOKIE:" + strconv.Itoa(cookie))
		} else {
			reply = []byte("STATUS:RECONNECT\nCOOKIE:" + strconv.Itoa(cookie))
		}
		log.Println(cookie)
		p.ttl = 7200
		p.flag = true
		activeDict[cookie] = p

		reply = append(reply, byte('\r'))
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}
		log.Println("Active one's left are :")
		for i, _ := range activeDict {
			log.Println(i)
		}

	case "LEAVE":
		// TODO: Only receive cookie
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ":")
			switch s[0] {
			case "HOSTNAME":
				p.hostName = s[1]
			case "PORT":
				p.port, _ = strconv.Atoi(s[1])
			case "COOKIE":
				cookie, _ = strconv.Atoi(s[1])
			}
		}
		p.ttl = 0
		p.flag = false
		delete(activeDict, cookie)
		log.Println("Active one's left are :")
		for i, _ := range activeDict {
			log.Println(i)
		}

	case "GETNODES":
		var reply []byte
		for _, value := range activeDict {
			// TODO: y strconv all the places .. change port to string
			reply = append(reply, []byte(value.hostName+":"+strconv.Itoa(value.port)+"\n")...)
		}
		reply = append(reply, byte('\r'))
		if _, err := conn.Write(reply); err != nil {
			log.Println(err)
		}
		log.Println("Replied with: ", string(reply))
	}
}
