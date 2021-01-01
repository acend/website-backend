package main

import (
	"backend/internal/config"
	"backend/internal/server"
)

func main() {
	config.Configure()
	server.Start()
}
