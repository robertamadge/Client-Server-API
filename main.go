package main

import (
	client "github.com/robertamadge/Client-Server-API/Client"
	server "github.com/robertamadge/Client-Server-API/Server"
)

func main() {
	go server.RunServer() // Run the server concurrently
	client.RunClient()
}
