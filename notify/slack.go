package notify

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Slack struct {
	webhookURL string
}

func New(webhookURL string) *Slack {
	return &Slack{webhookURL: webhookURL}
}

type payload struct {
	Text string `json:"text"`
}

func (s *Slack) Do (message string) error {
	p, err := json.Marshal(payload{Text: message})
	if err != nil {
		return err
	}
	resp, err := http.PostForm(s.webhookURL, url.Values{"payload": {string(p)}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
