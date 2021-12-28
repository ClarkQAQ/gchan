package gchan

import "errors"

type Message struct {
	ID   string `json:"i"` // 编号
	Chan string `json:"c"` // 通道
	Type int    `json:"t"` // 消息类型
	Body string `json:"b"` // 消息体
}

const (
	TypeRequest       int = 0x1A // 请求消息
	TypeZeroRequest   int = 0x1B // 无返回请求消息
	TypeResponse      int = 0x2A // 返回消息
	TypeErrorResponse int = 0x2B // 错误返回消息
)

var (
	ErrFuncCantBeNil        = errors.New("function cannot be nil")
	ErrChanNotFound         = errors.New("channel not found")
	ErrRequestTimeout       = errors.New("request timeout")
	ErrResponseIsNil        = errors.New("response is nil")
	ErrClientNotSupported   = errors.New("client type not supported")
	ErrClientNotFound       = errors.New("not found client")
	ErrRequestAlreadySent   = errors.New("request already sent")
	ErrRequestIsZeroRequest = errors.New("request is zero request")
)
