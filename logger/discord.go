package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Discord struct {
	client  *http.Client
	url     string
	msgChan chan []byte
	exit    chan struct{}
	closed  bool
}

func NewDiscord(webUrl string) *Discord {
	return &Discord{
		client:  http.DefaultClient,
		url:     webUrl,
		msgChan: make(chan []byte, 1000),
		exit:    make(chan struct{}, 1),
		closed:  false,
	}
}

func (d *Discord) Send(msg string) {
	if d.closed || len(d.msgChan) == 1000 {
		return
	}
	data, _ := json.Marshal(DisMsg{
		Content: fmt.Sprintf("%v  %v", time.Now().Format("2006-01-02 15:04:05"), msg),
	})
	d.msgChan <- data
}

func (d *Discord) Run() {
	for {
		select {
		case message, ok := <-d.msgChan:
			if ok {
				_, err := d.Call(message)
				if err != nil {
				}
			}
		case <-d.exit:
			return
		}
	}
}
func (d *Discord) Close() error {
	d.closed = true
	close(d.exit)
	return nil
}

func (d *Discord) Call(message []byte) (string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, d.url, bytes.NewBuffer(message))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := d.client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("send discord error %v %s", resp.StatusCode, all)
	}
	return string(all), nil
}

type DisMsg struct {
	Content string `json:"content"`
}

func newDisMsg(msg string) []byte {
	m := fmt.Sprintf(`{"content":"%v  %v"}`, time.Now().Format("2006-01-02 15:04:05"), msg)
	return []byte(m)
}
