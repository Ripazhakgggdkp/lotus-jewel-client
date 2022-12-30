package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	client, err := newClient()
	if err != nil {
		log.Println("Could not start client connection. Make sure you are running Intiface Central.\n", err)
		exit()
	}

	log.Println("Connection OK")

	clientErr := make(chan error)
	go func(err chan error) {
		clientErr <- client.handle()
	}(clientErr)

	httpErr := listenHTTP(client)

	for {
		select {
		case err := <-clientErr:
			log.Println(err)
			exit()
		case err := <-httpErr:
			log.Println(err)
			exit()
		}
	}
}

func exit() {
	fmt.Println("Press enter key to exit")
	fmt.Scanln()
	os.Exit(1)
}
