package routes

import (
	"encoding/json"
	"github.com/fasthttp-contrib/websocket"
	"olivia/models"
	"strconv"
	"log"
)

type WSClient struct {
	Conn *websocket.Conn
	User *models.User
}

type WSMessage struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

var WSClients []WSClient

func WS(c *websocket.Conn) {
	defer c.Close()

	for {
		var message WSMessage
		mt, text, err := c.ReadMessage()

		// Handle disconnection
		if mt == websocket.CloseGoingAway || mt == -1 {
			removingKey := -1
			for i, client := range WSClients {
				if client.Conn == c {
					removingKey = i
				}
			}

			if removingKey != -1 {
				WSClients[removingKey] = WSClients[len(WSClients)-1]
				WSClients = WSClients[:len(WSClients)-1]

				//log.Println("Removing client...")
				//log.Printf("%+v", WSClients)
			}

			c.Close()
			break
		}

		if err != nil {
			break
		}

		json.Unmarshal(text, &message)

		if message.Type == "hello" {
			userId, _ := strconv.ParseInt(message.Data["userId"], 10, 64)

			var user models.User
			user.Id = userId
			user.GetById()

			WSClients = append(WSClients, WSClient{
				Conn: c,
				User: &user,
			})

			//log.Println("New client...")
			//log.Printf("%+v", WSClients)
		}
	}
}

func SendWSToUser(user models.User, data interface{}) bool {
	for _, client := range WSClients {
		if client.User.Id == user.Id {
			log.Printf("Sending WS message to %s", user.FirstName)
			client.Conn.WriteJSON(data)
		}
	}

	return true
}
