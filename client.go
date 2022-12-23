package main

import (
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

func newClient() (*Client, error) {
	socketUrl := "ws://localhost:12345"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		LastActivated: time.Now(),
		conn:          conn,
		connected:     false,
	}, nil
}

// Keeps the websocket connection alive
func (g *Client) ping() <-chan error {
	error := make(chan error)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			if err := g.send(websocket.PingMessage, []byte{}); err != nil {
				error <- err
			}
		}
	}()
	return error
}

func (g *Client) connect(ID int) error {
	g.connected = true
	return g.send(websocket.TextMessage, connect(ID))
}

func (g *Client) vibrate(strength float64) error {
	g.LastActivated = time.Now()
	return g.send(websocket.TextMessage, vibrate(1, clamp(strength, 0, 1)))
}

// Makes sure device activation only runs for a set amount of time
func (g *Client) autoStop(duration time.Duration, ID int) <-chan error {
	error := make(chan error)
	go func(g *Client) {
		for {
			if time.Since(g.LastActivated) > duration && g.connected {
				err := g.send(websocket.TextMessage, stop(ID))
				if err != nil {
					error <- err
				}
				g.lock.Lock()
				g.LastActivated = time.Now()
				g.lock.Unlock()
			}
		}
	}(g)

	return error
}

func (g *Client) stop() error {
	return g.send(websocket.TextMessage, stop(1))
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

func (g *Client) send(messageType int, data []byte) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.conn.WriteMessage(messageType, data)
}
