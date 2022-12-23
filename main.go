package main

import (
	"time"
)

var duration = 1 * time.Second
var ID = 1

func main() {

	client := newClient()

	client.ping()
	client.autoStop(duration, ID)

	listenHTTP(client)

}
