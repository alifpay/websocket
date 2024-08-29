package main

import (
	"log"
	"net/http"

	"github.com/alifpay/websocket/app"
)

func main() {
	http.HandleFunc("/ws", app.Mw(app.WebSocket))
	
	//websocket broadcast
	go app.Broadcast()

	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
