package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/checker"
	"net-cat/internal/tcpserver"
)

var (
	MESSAGE = `[USAGE]: ./TCPChat $port`
	PORT    = `8989`
	HOST    = `localhost`
	TYPE    = `tcp`
)

func main() {
	parameters := os.Args[1:]
	if len(parameters) > 0 && !checker.Valid(parameters) {
		log.Fatal(MESSAGE)
		os.Exit(1)
	}

	if len(parameters) == 1 {
		PORT = parameters[0]
	}
	log.SetFlags(0)
	server := tcpserver.NewServer()
	if err := server.Listen(TYPE, fmt.Sprintf("%v:%v", HOST, PORT)); err != nil {
		log.Fatal("Error starting server:", err)
	}

	defer server.Close()
	server.Start()
}
