package main

import (
	"fmt"
	"gchan"
	"time"
)

func main() {
	gh := gchan.New()

	gh.Reader(func(b []byte) error {
		fmt.Printf("[MESSAGE] 时间: %s 内容: %v\n",
			time.Now().Format("2006-01-02 15:04:05"), string(b))
		go gh.Writer(b)
		return nil
	})

	gh.Get("hello", func(c *gchan.Context) {
		c.String("hello " + c.GetString())
	})

	for i := 0; i < 20; i++ {
		go func() {
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
		}()
	}
	fmt.Scanln()
}
