package logger

import (
	"testing"
	"time"
)

func TestDiscord(t *testing.T) {
	//url := "https://discord.com/api/webhooks/1260480229223436288/JMMCTRMXnV2miHmGJaxMiGGK2EeemzQZtX6JpcqEvemcKmhRPCu9fv6qmiomeYhLVOic"
	//url := "https://discord.com/api/webhooks/1260514595974676490/I8XADjbJyPQr2MP0qzOH1W3yawsfnhzx_xcA6Joo1dShQEv6UdJukuc6sKmDRsGPRXwJ"
	url := "https://discord.com/api/webhooks/1260515536324792362/FBuJAgPaXu3q4bXGVLswxvNt3m8EyzCkx195Z1gtZeAPpZAW8GztCDTbEFrrOYvtjq7_"
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

func TestDemo(t *testing.T) {

}
