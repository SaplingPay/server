package handlers

import (
	"errors"
	"fmt"
	"github.com/SaplingPay/server/models"
	"github.com/SaplingPay/server/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"SaplingPay"},
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins, be cautious with this in production
		// TODO - make this more secure
		return true
	},
}

func GetTableWSSession(c *gin.Context) {
	roomID := c.Param("tableId")
	bearerToken := c.Request.Header.Get("Authorization")

	fmt.Printf("Room ID: %s\n %s", roomID, bearerToken)

	hub := ws.GetHub(roomID)

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	defer func() {
		hub.RemoveClient(conn)
	}()

	hub.AddClient(conn)

	handleMessages(hub, conn)
}

func handleMessages(hub *ws.TableHub, conn *websocket.Conn) {
	for {
		// Read message from browser
		var message models.TableMessage
		err := conn.ReadJSON(&message)

		if !errors.Is(err, nil) {
			log.Printf("error occurred: %v", err)
			hub.RemoveClient(conn)
			break
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), message.Message)

		hub.Broadcast <- ws.LocalMessage{Message: message, Conn: conn}
	}
}
