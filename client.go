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
	Type   int        `json:"type"`
	Body   string     `json:"body,omitempty"`
	Slider SliderMove `json:"slider,omitempty"`
}

type SliderMove struct {
	SliderID     int     `json:"id, omitempty"`
	PercentValue float32 `json:"percent, omitempty"`
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

		//bytes := []byte(string(p))
		log.Printf(string(p))

		//Check for slider
		var sliderMove SliderMove
		sliderError := json.Unmarshal(p, &sliderMove)
		if sliderError != nil {
			log.Println("no json")
		}

		// { "SliderID": 1, "PercentValue": 0.2}
		// {"id": 1, "percent": 0.2}
		//sliderMoveStatic := SliderMove{SliderID: 1, PercentValue: 0}
		message := Message{Type: messageType, Body: string(p), Slider: sliderMove}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)

	}
}
