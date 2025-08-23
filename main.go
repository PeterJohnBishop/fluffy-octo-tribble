package main

import (
	"fluffy-coto-tribble/server"
	"fluffy-coto-tribble/server/authentication"
	"log"
)

func main() {
	log.Println("Hello, Fluffy Octo Tribble!")
	authentication.InitAuth()
	server.InitServer()
}
