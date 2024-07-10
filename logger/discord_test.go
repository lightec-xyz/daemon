package logger

import (
	"testing"
	"time"
)

func TestDiscord(t *testing.T) {
	discord := NewDiscord("https://discord.com/api/webhooks/1260480229223436288/JMMCTRMXnV2miHmGJaxMiGGK2EeemzQZtX6JpcqEvemcKmhRPCu9fv6qmiomeYhLVOic")
	_, err := discord.Call(newDisMsg("running ...."))
	if err != nil {
		t.Error(err)
	}
	go discord.Run()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			discord.Send("hello")
		}
	}()
	ch := make(chan struct{}, 1)
	<-ch

}

func TestDemo(t *testing.T) {

}
