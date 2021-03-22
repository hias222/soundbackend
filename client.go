package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body,omitempty"`
}

type SoundMessage struct {
	Type    int `json:"type"`
	Message SliderMessage
}
type SliderMessage struct {
	SliderID     int     `json:"id,omitempty"`
	PercentValue float32 `json:"percent,omitempty"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var jsonmessage bool

		jsonmessage = true

		//bytes := []byte(string(p))
		log.Printf(string(p))

		//Check for slider
		var sliderMessage SliderMessage
		sliderError := json.Unmarshal(p, &sliderMessage)
		if sliderError != nil {
			jsonmessage = false
			log.Println("no json")
		}

		if jsonmessage {
			soundmessage := SoundMessage{Type: 2, Message: sliderMessage}
			c.Pool.Soundcast <- soundmessage
			fmt.Printf("New Sound Message : %+v\n", soundmessage)
		}

		// {"id": 1, "percent": 0.2}
		//sliderMoveStatic := SliderMove{SliderID: 1, PercentValue: 0}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)

	}
}
