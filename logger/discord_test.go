package logger

import (
	"testing"
	"time"
)

func TestDiscord(t *testing.T) {
	url := ""
	discord := NewDiscord(url)
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
