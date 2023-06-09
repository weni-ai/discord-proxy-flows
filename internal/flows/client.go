package flows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Message struct {
	Text        string   `json:"text,omitempty"`
	From        string   `json:"from,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) SendReceiveMessage(message Message) error {
	url := fmt.Sprintf("%s/receive", c.BaseURL)
	return c.sendMessageForm(url, message.Text, message.From, message.Attachments)
}

func (c *Client) SendSentMessage(message FlowsMessage) error {
	url := fmt.Sprintf("%s/sent", c.BaseURL)
	return c.sendMessage(url, message)
}

func (c *Client) sendMessage(url string, message interface{}) error {
	jsonStr, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) sendMessageForm(urlDestination, text, from string, attachments []string) error {

	formData := url.Values{
		"text": {text},
		"from": {from},
	}

	if attachments != nil && len(attachments) > 0 {
		formData = url.Values{
			"text":        {text},
			"from":        {from},
			"attachments": attachments,
		}
	}

	resp, err := http.PostForm(urlDestination, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type FlowsMessage struct {
	ID          string   `json:"id,omitempty"`
	Text        string   `json:"text,omitempty"`
	To          string   `json:"to,omitempty"`
	Channel     string   `json:"channel,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}
