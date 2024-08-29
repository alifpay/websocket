package app

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type CtxUserKey struct{}

var wsClients = make(map[*websocket.Conn]string)
var broadcast = make(chan Message)
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Action  string `json:"action"`
	User    string `json:"user"`
	Payload any    `json:"payload"`
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	userName := r.Context().Value(CtxUserKey{}).(string)
	if len(userName) == 0 {
		log.Println("User name not found in context")
		return
	}
	//todo get user id from db by user name
	wsClients[ws] = userName
	for {
		var req Message
		err := ws.ReadJSON(&req)
		if err != nil {
			log.Println("ws.ReadJSON: ", err)
			delete(wsClients, ws)
			break
		}
		req.User = userName
		resp, err := ProcessMessage(req)
		if err != nil {
			log.Println("ProcessMessage: ", err)
			continue
		}
		if resp.Payload != nil {
			broadcast <- resp
		}
	}
}

// broadcast messages to all clients or to specific user
func Broadcast() {
	for {
		msg := <-broadcast
		for client, user := range wsClients {
			if len(msg.User) > 0 && user != msg.User {
				continue
			}
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("client.WriteJSON: ", err)
				client.Close()
				delete(wsClients, client)
			}
		}
	}
}

func Mw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tokenTEststsdfs_sdf_dfsdf_sdsfdsdf234243" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("Authorization successful")
		next.ServeHTTP(w, r)
	}
}

func ProcessMessage(req Message) (resp Message, err error) {
	switch req.Action {
	case "ping":
		resp.Action = "pong"
	case "pong":
		resp.Action = "ping"
	}
	return
}
