package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	LastActivated time.Time
	conn          *websocket.Conn
	lock          sync.Mutex
	connected     bool
}

func newClient() *Client {
	socketUrl := "ws://localhost:12345"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Could not connect to websocket server. Check if Intiface Central is installed and running.\n", err)
	}
	return &Client{
		LastActivated: time.Now(),
		conn:          conn,
		connected:     false,
	}
}

// Keeps the websocket connection alive
func (g *Client) ping() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			if err := g.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}()
}

func (g *Client) connect(ID int) {
	g.conn.WriteMessage(websocket.TextMessage, connect(ID))
	g.connected = true
}

func (g *Client) vibrate(strength float64) {
	g.LastActivated = time.Now()
	g.conn.WriteMessage(websocket.TextMessage, vibrate(1, clamp(strength, 0, 1)))
}

// Makes sure device activation only runs for a set amount of time
func (g *Client) autoStop(duration time.Duration, ID int) {
	go func(g *Client) {
		for {
			if time.Since(g.LastActivated) > duration && g.connected {
				g.conn.WriteMessage(websocket.TextMessage, stop(ID))
				g.lock.Lock()
				g.LastActivated = time.Now()
				g.lock.Unlock()
			}
		}
	}(g)

}

func (g *Client) stop() {
	g.conn.WriteMessage(websocket.TextMessage, stop(1))
}

func clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}
