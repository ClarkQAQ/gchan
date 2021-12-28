package gchan

import (
	"encoding/json"
)

func (g *Gchan) handle(b []byte) error {
	if len(b) <= 0 || b[0] != '{' || b[len(b)-1] != '}' {
		return nil
	}

	msg := Message{}
	if e := json.Unmarshal(b, &msg); e != nil {
		return e
	}

	switch msg.Type {
	case TypeRequest:
		return g.handleRequest(&msg)
	case TypeZeroRequest:
		return g.handleRequest(&msg)
	case TypeResponse:
		return g.handleResponse(&msg)
	case TypeErrorResponse:
		return g.handleResponse(&msg)
	}

	return nil
}

func (g *Gchan) handleRequest(msg *Message) error {
	g.mLock.RLock()
	h, ok := g.m[msg.Chan]
	g.mLock.RUnlock()

	if !ok || h == nil {
		// 类似 HTTP 的404 Not Found
		return g.sendMsg(msg.resMsg(TypeErrorResponse,
			ErrChanNotFound.Error()))
	}

	// 创建一个新的上下文
	// 然后剩下的事情交给处理函数去做
	h(g.newContext(msg))
	return nil
}

func (g *Gchan) handleResponse(msg *Message) error {
	if v, ok := g.c.Load(msg.ID); ok {

		// 把整个消息结构写进channel
		if ch, ok := v.(chan *Message); ok {
			ch <- msg
			return nil
		}

		return ErrClientNotSupported
	}

	return ErrClientNotFound
}
