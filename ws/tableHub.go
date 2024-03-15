package ws

import (
	"errors"
	"github.com/SaplingPay/server/models"
	"github.com/gorilla/websocket"
	"log"
)

type LocalMessage struct {
	Message models.TableMessage
	Conn    *websocket.Conn
}

type TableHub struct {
	roomId    string
	Clients   map[*websocket.Conn]bool
	Broadcast chan LocalMessage
}

func NewHub(roomId string) *TableHub {
	return &TableHub{
		roomId:    roomId,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan LocalMessage),
	}
}

func (h *TableHub) Run() {
	for {
		select {
		case message := <-h.Broadcast:
			for client := range h.Clients {
				// Don't forward the message to the sender
				if client == message.Conn {
					continue
				}

				if err := client.WriteJSON(message.Message); !errors.Is(err, nil) {
					log.Printf("error occurred: %v", err)
				}
			}
		}
	}
}

func (h *TableHub) Close() {
	for client := range h.Clients {
		client.Close()
	}
}

func (h *TableHub) AddClient(client *websocket.Conn) {
	h.Clients[client] = true
}

func (h *TableHub) RemoveClient(client *websocket.Conn) {
	delete(h.Clients, client)
	client.Close()
}

// Global Meta-Hub for all hubs
var wsHubs = make(map[string]*TableHub)

func GetHub(roomID string) *TableHub {
	if wsHubs[roomID] == nil {
		wsHubs[roomID] = NewHub(roomID)
		go wsHubs[roomID].Run()
	}
	return wsHubs[roomID]
}

func RemoveHub(roomID string) {
	if wsHubs[roomID] != nil {
		wsHubs[roomID].Close()
		delete(wsHubs, roomID)
	}
}
