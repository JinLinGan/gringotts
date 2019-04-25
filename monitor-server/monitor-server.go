package main

import (
	"log"

	"github.com/jinlingan/gringotts/server"
)

const (
	serverID = "99"
)

func main() {
	address := ":7777"
	serverInst, err := server.NewServer(address, serverID)
	if err != nil {
		log.Fatalf("can not create new server in port %s : %s", address, err)
	}
	log.Fatal(serverInst.Serve())
}
