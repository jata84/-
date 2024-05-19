package main

import (
	"goTask/server"
)

func main() {
	server := server.NewServer()
	server.Run()
}
