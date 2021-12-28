package gchan

import (
	"encoding/json"
)

func newMsg(id, ch string, tp int, body string) *Message {
	return &Message{
		ID:   id,
		Chan: ch,
		Type: tp,
		Body: body,
	}
}

func (g *Gchan) sendMsg(msg *Message) error {
	b, e := json.Marshal(msg)
	if e != nil {
		return e
	}

	return g.reader(b)
}

func (m *Message) resMsg(tp int, body string) *Message {
	return newMsg(m.ID, m.Chan, tp, body)
}
