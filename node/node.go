package main

import (
	"flag"
	"hrank/ptp/peers"
	"log"
)

func main() {
	// Get the cli arguments and populate peer details
	hostName := flag.String("hostname", "localhost", "host name of peer")
	port := flag.Int("port", 85000, "Port Num of peer")
	data := flag.String("data", "", "data associated with peer")

	// Do this on actual file send
	//if err != nil {
	//	log.Fatal("File corrupt")
	//}

	// Check reconnect or new connection
	p, peerData := peers.NewPeer(*port, *hostName, *data)
	log.Println("Peer Registering with RS")
	if err = p.Register(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Peer %s registered at port %d\n", p.hostName, p.port)
}
