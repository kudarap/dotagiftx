package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client represents discord client.
type Client struct {
	webhookURL string
}

func New(webhookURL string) *Client {
	return &Client{webhookURL}
}

func (c *Client) PostWebhook(username, content string) error {
	if username == "" || content == "" {
		return fmt.Errorf("username and content required")
	}

	payload := struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}{
		username,
		content,
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}
	resp, err := http.Post(c.webhookURL, "application/json", b)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
