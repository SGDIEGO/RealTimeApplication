package main

import (
	"log"
	"net"
	"sync"

	"github.com/SGDIEGO/RealTimeApp/pkg/server"
)

func main() {

	log.Println("Connecting to server")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	var connMap = &sync.Map{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := server.Client(conn)
		client.ListenMssge(connMap)
	}
}
