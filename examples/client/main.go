package main

import (
	"fmt"
	"gchan"
	"os"

	"golang.org/x/net/websocket"
)

const (
	host = "127.0.0.1:8080"
)

func main() {
	fmt.Println("client start")

	ws, e := websocket.Dial("ws://"+host+"/", "", "http://"+host+"/")
	if e != nil {
		fmt.Printf("dial error: %v\n", e)
		os.Exit(1)
	}
	ws.PayloadType = websocket.TextFrame

	fmt.Println("client connect success")

	gh := gchan.New()

	gh.Reader(func(b []byte) error {
		return websocket.Message.Send(ws, b)
	})

	go func() {
		for {
			var data []byte
			if e := websocket.Message.Receive(ws, &data); e != nil {
				ws.Close()
				os.Exit(1)
			}

			gh.Writer(data)
		}
	}()

	for {
		ret := gh.Set("hello")
		ret = ret.String(fmt.Sprintf("uuid: %s", ret.GetID()))
		res := ret.End()
		if e := res.Error(); e != nil {
			fmt.Printf("request error: %v\n", e)
			continue
		}

		fmt.Printf("response string: %s\n", res.String())
	}
}
