package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Send struct {
	websocketMessageType int
	messageType          Message
	data                 []byte
}

type Message string

const (
	UnknownMessage   Message = ""
	VibrateMessage   Message = "vibrate"
	PingMessage      Message = "ping"
	StopMessage      Message = "stop"
	ConnectMessage   Message = "connect"
	EnumerateMessage Message = "enumerate"
)

type SendError struct {
	messageType Message
	err         error
}

type Client struct {
	LastActivated time.Time
	conn          *websocket.Conn
	connected     bool
	devices       []Devices
	sendMsg       chan Send
	sendError     chan SendError
	sync.Mutex
}

var duration = 1 * time.Second
var ID = 1

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

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
		sendError:     make(chan SendError),
		sendMsg:       make(chan Send),
		devices:       nil,
	}, nil
}

// Keeps the websocket connection alive
func (g *Client) ping() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		g.send(websocket.PingMessage, PingMessage, []byte{})
	}
}

// Wait for device enumeration from Intiface
func (c *Client) getDevices() {
	c.Lock()
	c.devices = nil
	c.connected = false
	c.Unlock()

	response := []DeviceListResponse{}
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.sendError <- SendError{messageType: EnumerateMessage}
			}
			break
		}

		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Fatal(err)
		}

		devices := response[0].DeviceList

		if (devices.ID) == 1 {
			log.Println("Connected devices:", devices.Devices)
			c.Lock()
			c.devices = response[0].DeviceList.Devices
			c.connected = true
			c.Unlock()
			break
		}
	}
}

func (g *Client) connect(ID int) error {
	g.send(websocket.TextMessage, ConnectMessage, connect(ID))
	log.Println("Connecting game to Intiface Central")

	g.send(websocket.TextMessage, ConnectMessage, enumerate(ID))
	log.Println("Requesting connected devices")

	go g.getDevices()

	return nil

}

func (g *Client) vibrate(strength float64) {
	g.Lock()
	g.LastActivated = time.Now()
	indexes := g.devices
	g.Unlock()

	for _, device := range indexes {
		g.send(websocket.TextMessage, VibrateMessage, vibrate(1, device.DeviceIndex, clamp(strength, 0, 1)))
	}
}

// Makes sure device activation only runs for a set amount of time
func (g *Client) autoStop(duration time.Duration, ID int) {
	for {
		g.Lock()
		if time.Since(g.LastActivated) > duration && g.connected {
			for _, device := range g.devices {
				g.send(websocket.TextMessage, StopMessage, stop(device.DeviceIndex))
			}

			g.LastActivated = time.Now()
		}
		g.Unlock()
	}
}

func (g *Client) stop() {
	g.Lock()
	indexes := g.devices
	g.Unlock()

	for _, device := range indexes {
		g.send(websocket.TextMessage, StopMessage, stop(device.DeviceIndex))
	}

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

func (g *Client) send(websocketMessageType int, messageType Message, data []byte) {
	g.sendMsg <- Send{websocketMessageType: websocketMessageType, messageType: messageType, data: data}
}

func (g *Client) handle() error {
	g.conn.SetReadLimit(maxMessageSize)
	g.conn.SetReadDeadline(time.Now().Add(pongWait))
	g.conn.SetPongHandler(func(string) error { g.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go g.ping()
	go g.autoStop(duration, ID)

	for {
		select {
		case err := <-g.sendError:
			fmt.Println(err)
			return err.err
		case msg := <-g.sendMsg:
			if err := g.conn.WriteMessage(msg.websocketMessageType, msg.data); err != nil {
				return err
			}
		}
	}
}
