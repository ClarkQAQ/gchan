#### [Gchan] golang 轻量级双向同步请求

> 苦于写 websocket 的时候一堆一堆的 switch/map 和 难以处理的同步请求,要是 websocket 开发和 http 一样就好了, 这就是 Gchan

> 我没有绑定任何底层, 只需要一个 bytes 的 io 接口就行 (后面打算尝试封装wasm).....x/websocket也只是演示调用


##### 描述

1.由于 "轻量化" 的关系, 把本来写好的树结构路由给砍了.....现在的路由只能完全匹配字符串不能通配之类的了

2.一切实现都是仿照我之前的 web 框架实现的,最大化模拟 http 请求的调用

3.框架传输用的是 json, 当初想好要跨平台.....后面再加上多个codec

4.用在 tcp 或者其他会"粘包" 的环境要自己处理包分割....

5.框架内部使用的是 goroutine, 可以自己控制 `Writer` 进而并发数量, 因为框架内部没有线程池, 所以框架内部的并发数量是控制不了的

6.请求的 `Timeout` 设置为0的时候, 服务端将不会返回, 客户端也不会等待....


##### 示例

1.服务端:


```golang

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


```

2.客户端:

```golang

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


```