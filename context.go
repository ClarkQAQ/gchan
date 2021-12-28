package gchan

import (
	"encoding/json"
)

type Handler func(c *Context)

type Context struct {
	m *Message
	g *Gchan
}

func (g *Gchan) newContext(m *Message) *Context {
	return &Context{
		m: m,
		g: g,
	}
}

func (c *Context) Chan() string {
	return c.m.Chan
}

func (c *Context) ID() string {
	return c.m.ID
}

func (c *Context) GetString() string {
	return c.m.Body
}

func (c *Context) GetJSON(v interface{}) error {
	return json.Unmarshal([]byte(c.m.Body), v)
}

func (c *Context) String(s string) error {
	if c.m.Type == TypeZeroRequest {
		return ErrRequestIsZeroRequest
	}

	return c.g.sendMsg(c.m.resMsg(TypeResponse, s))
}

func (c *Context) JSON(v interface{}) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}

	return c.String(string(b))
}
