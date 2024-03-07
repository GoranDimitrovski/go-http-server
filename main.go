package main

import (
	"fmt"
	"log"
	"net/http"
	"simplesurance/config"
	"simplesurance/timestamp"
)

func main() {
	timestamp.Init(config.ServerConfig.Filename)
	http.HandleFunc(config.ServerConfig.Route, timestamp.Handler)

	serverAddr := fmt.Sprintf(":%s", config.ServerConfig.Port)
	log.Printf("Starting server at %s%s\n", config.ServerConfig.Address, config.ServerConfig.Route)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
