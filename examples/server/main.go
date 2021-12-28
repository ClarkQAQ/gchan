package main

import (
	"fmt"
	"gchan"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func handler(ws *websocket.Conn) {
	ws.PayloadType = websocket.TextFrame
	gh := gchan.New()

	gh.Reader(func(b []byte) error {
		return websocket.Message.Send(ws, b)
	})

	gh.Get("hello", func(c *gchan.Context) {
		c.String("hello " + c.GetString())
	})

	for {
		var data []byte
		if e := websocket.Message.Receive(ws, &data); e != nil {
			ws.Close()
			return
		}

		gh.Writer(data)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(handler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("server start")
	if e := server.ListenAndServe(); e != nil {
		fmt.Printf("listen error: %v\n", e)
		os.Exit(1)
	}
}
