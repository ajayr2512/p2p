package main

import (
	"flag"
	"fmt"
	"hrank/ptp/peers"
	"log"
)

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
	log.Println("From Node : ")
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
			break loop
		}

	}
}
