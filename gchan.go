package gchan

import (
	"runtime"
	"sync"
)

type Gchan struct {
	c     *sync.Map          // 客户端表
	m     map[string]Handler // 路由表
	mLock *sync.RWMutex      // 路由表锁

	reader func([]byte) error // 输出回调
}

func New() *Gchan {
	return &Gchan{reader: func(b []byte) error {
		return nil
	},
		c:     &sync.Map{},
		m:     make(map[string]Handler),
		mLock: &sync.RWMutex{},
	}
}

// 输出回调
func (g *Gchan) Reader(w func([]byte) error) error {
	if w == nil {
		return ErrFuncCantBeNil
	}

	g.reader = func(b []byte) error {
		return w(b)
	}
	return nil
}

// 输入数据
func (g *Gchan) Writer(b []byte) error {
	return g.handle(b)
}

func (g *Gchan) Close() {
	g.mLock.Lock()
	defer g.mLock.Unlock()

	// 删除路由
	for k := range g.m {
		delete(g.m, k)
	}

	// 清空待返回的客户端channel
	g.c.Range(func(k, v interface{}) bool {
		g.c.Delete(k)
		return true
	})

	// 时停一下
	runtime.GC()
}

// [服务端] 设置监听获取数据
func (g *Gchan) Get(ch string, h Handler) {
	g.mLock.Lock()
	defer g.mLock.Unlock()

	// 如果给的是空handler，则删除该路由
	if h == nil {
		delete(g.m, ch)
		return
	}

	// 加入路由
	g.m[ch] = h
}

// [客户端] 发送数据
func (g *Gchan) Set(ch string) *Request {
	return g.newClient(ch)
}
