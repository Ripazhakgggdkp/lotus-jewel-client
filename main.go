package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var duration = 1 * time.Second
var ID = 1

func main() {

	client, err := newClient()
	if err != nil {
		log.Println("Could not start client connection. Make sure you are running Intiface Central.\n", err)
		exit()
	}

	log.Println("Connection OK")

	pingErr := client.ping()
	go func(err chan error) {
		log.Println("Could not keep connection alive, exiting\n", <-err)
		exit()
	}(pingErr)

	autoStopErr := client.autoStop(duration, ID)
	go func(err chan error) {
		log.Println("Could not stop device, exiting\n", <-err)
		exit()
	}(autoStopErr)

	httpErr := listenHTTP(client)
	if <-httpErr != nil {
		log.Println("Server error\n", err)
	}
}

func exit() {
	fmt.Println("Press enter key to exit")
	fmt.Scanln()
	os.Exit(1)
}
