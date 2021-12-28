package gchan

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	timeout time.Duration // 计划超时时间
	err     error         // 流程错误
	isend   bool          // 防止重复发送
	m       *Message      // 发送消息结构体
	g       *Gchan
}

type Response struct {
	err error    // 返回错误
	req *Request // 发送请求存档
	m   *Message // 接收到的消息
}

func (g *Gchan) newClient(ch string) *Request {
	return &Request{
		// 默认超时时间为 6 秒
		// 6 秒后收不到回复
		// 则认为请求超时并释放等待的请求
		timeout: time.Second * 6,
		// 创建一个新的请求
		// Id 用uuid生成,也可以用雪花,但是我懒得写了
		m: newMsg(uuid.New().String(), ch, TypeRequest, ""),
		g: g,
	}
}

// 超时时间
// 例子: ret.Timeout(time.Second * 10)
func (ret *Request) Timeout(timeout time.Duration) *Request {
	ret.timeout = timeout

	// 当超时定义为 0 时,则发送无返回请求
	if ret.timeout <= 0 {
		ret.m.Type = TypeZeroRequest
	}
	return ret
}

// 获取请求编号
// 不可自定义因为这样可能会发生严重的错误
func (ret *Request) GetID() string {
	return ret.m.ID
}

// 写入文本数据
func (ret *Request) String(s string) *Request {
	ret.m.Body = s
	return ret
}

// 写入JSON数据
// 此处无视了序列化结构体的错误
func (ret *Request) JSON(v interface{}) *Request {
	b, e := json.Marshal(v)
	ret.err = e
	ret.m.Body = string(b)
	return ret
}

// 发送请求
// 这一步才是真正的发送请求过去
func (ret *Request) End() *Response {
	res := &Response{
		req: ret,
	}

	if ret.err != nil { // 如果有错误,那就不发了
		return res
	} else if ret.isend { // 如果已经发送过了,就不再发送
		ret.err = ErrRequestAlreadySent
		return res
	}

	// 定义一个 message channel 用于接收回复
	msgChan := make(chan *Message, 1)
	ret.g.c.Store(ret.m.ID, msgChan)

	defer func() {
		// 关闭时清理垃圾
		close(msgChan)
		ret.g.c.Delete(ret.m.ID)
	}()

	// 发送数据
	ret.isend = true
	if res.err = ret.g.sendMsg(ret.m); res.err != nil {
		return res
	}

	// 如果设置了超时时间为零则不待返回结果直接退出
	if ret.timeout == 0 {
		return res
	}

	select {
	case res.m = <-msgChan: // 接收channel
		return res
	case <-time.After(ret.timeout): // 超时
		res.err = ErrRequestTimeout
		return res
	}
}

// 获取字节数据
func (ret *Response) Error() error {
	if ret.err != nil {
		return ret.err
	}

	if ret.req.timeout == 0 {
		return nil
	}

	if ret.m == nil {
		return ErrResponseIsNil
	}

	if ret.m.Type == TypeErrorResponse {
		return errors.New(ret.m.Body)
	}

	return nil
}

// 获取文本数据
func (ret *Response) String() string {
	if ret.m != nil {
		return ret.m.Body
	}
	return ""
}

// 获取结构数据
func (ret *Response) JSON(v interface{}) error {
	if ret.m != nil {
		return json.Unmarshal([]byte(ret.m.Body), v)
	}

	return ErrResponseIsNil
}
