package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

var rooms = []string{"channel1", "channel2", "channel3 "}

type Message struct {
	Channel string `json:"channel"`
	Msg     string `json:"msg"`
	Userid  int    `json:"userId"`
}

func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg Message) {
			log.Println("channel:", msg.Channel, "msg:", msg.Msg)
			so.Join(msg.Channel)
			so.Emit("chat message", msg)
			so.BroadcastTo(msg.Channel, "chat message", msg)
		})

		so.On("error", func() {
			log.Println("so error")
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
