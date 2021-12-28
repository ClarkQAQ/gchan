package gchan_test

import (
	"gchan"
	"testing"
)

func TestReader(t *testing.T) {
	gh := gchan.New()

	if gh.Reader(func(b []byte) error {
		gh.Writer(b)
		return nil
	}) != nil {
		t.Error("Can't Set Reader")
		t.Fail()
	}
}

func TestRecv(t *testing.T) {
	gh := gchan.New()

	gh.Get("test", func(c *gchan.Context) {
		c.String("Test Recv")
	})
}

func TestClose(t *testing.T) {
	gchan.New().Close()
}

func TestSend(t *testing.T) {
	gh := gchan.New()

	if gh.Reader(func(b []byte) error {
		gh.Writer(b)
		return nil
	}) != nil {
		t.Error("Can't Set Reader")
		t.Fail()
	}

	gh.Get("test", func(c *gchan.Context) {
		t.Log(c.GetString())
		c.String("Test Recv")
	})

	res := gh.Set("test").String("Test Send").End()
	if res.Error() != nil {
		t.Errorf("Send Error: %s", res.Error())
		t.Fail()
	}

	if res.String() != "Test Recv" {
		t.Errorf("Send Error: %s", res.String())
		t.Fail()
	}
}

func BenchmarkSend(b *testing.B) {
	gh := gchan.New()

	if gh.Reader(func(b []byte) error {
		gh.Writer(b)
		return nil
	}) != nil {
		b.Error("Can't Set Reader")
		b.Fail()
	}

	gh.Get("test", func(c *gchan.Context) {
		c.String("Test Recv")
	})

	for i := 0; i < b.N; i++ {
		gh.Set("test").String("Test Send").End()
	}
}
