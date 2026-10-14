// signaling.go
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			_ = client.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("Signaling Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
