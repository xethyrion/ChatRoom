package main

import (
	"ChatRoom/server"
	"fmt"
)

func main() {
	fmt.Println("A Simple ChatRoom")
	var Server server.XServer
	Server.SetPort(4545)
	Server.Start()
}
