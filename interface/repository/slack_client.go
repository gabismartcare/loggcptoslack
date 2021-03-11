package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const DefaultSlackTimeout = 5 * time.Second

func NewSlackClient(username, channel, webhook string) SlackClient {
	return SlackClient{
		UserName: username,
		Channel:  channel,
		TimeOut:  DefaultSlackTimeout,
		Url:      webhook,
	}
}

type FullSlackMessage struct {
	SlackMessage
	Username string `json:"username,omitempty"`
	Channel  string `json:"channel,omitempty"`
}
type SlackClient struct {
	UserName string
	Channel  string
	TimeOut  time.Duration
	Url      string
}

func (sc SlackClient) SendSlackMessage(s SlackMessage) error {
	slackRequest := FullSlackMessage{
		SlackMessage: s,
		Username:     sc.UserName,
		Channel:      sc.Channel,
	}
	slackBody, _ := json.Marshal(slackRequest)
	req, err := http.NewRequest(http.MethodPost, sc.Url, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.TimeOut == 0 {
		sc.TimeOut = DefaultSlackTimeout
	}
	client := &http.Client{Timeout: sc.TimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
