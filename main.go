package main

import (
	"fluffy-coto-tribble/server"
	"fluffy-coto-tribble/server/authentication"
	"log"
)

func main() {
	log.Println("Starting services...")
	authentication.InitAuth()
	server.InitServer()
}
