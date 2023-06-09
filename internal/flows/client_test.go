package flows

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_SendReceiveMessage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/receive" {
			t.Errorf("Expected route /receive, but got %s", r.URL.Path)
		}

		err := r.ParseForm()
		if err != nil {
			t.Error(err)
		}

		var receivedMessage Message

		receivedMessage.Text = getFormField(r.Form, "text")
		receivedMessage.From = getFormField(r.Form, "from")

		if receivedMessage.Text != "Hello" {
			t.Errorf("Expected message 'Hello', but got %s", receivedMessage.Text)
		}

		if receivedMessage.From != "112233445566778899" {
			t.Errorf("Expected sender '112233445566778899', but got %s", receivedMessage.From)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()

	client := NewClient(testServer.URL)

	message := Message{
		Text: "Hello",
		From: "112233445566778899",
	}

	err := client.SendReceiveMessage(message)
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
	}
}

func getFormField(form url.Values, name string) string {
	if name != "" {
		values, found := form[name]
		if found {
			return values[0]
		}
	}
	return ""
}

func TestClient_SendSentMessage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sent" {
			t.Errorf("Expected route /sent, but got %s", r.URL.Path)
		}

		var receivedMessage FlowsMessage
		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		if err != nil {
			t.Errorf("Error decoding request: %s", err.Error())
		}

		if receivedMessage.Text != "Hello" {
			t.Errorf("Expected message 'Hello', but got %s", receivedMessage.Text)
		}

		if receivedMessage.To != "112233445566778899" {
			t.Errorf("Expected sender '112233445566778899', but got %s", receivedMessage.To)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()

	client := NewClient(testServer.URL)

	flowsMessage := FlowsMessage{
		ID:      "123",
		Text:    "Hello",
		To:      "112233445566778899",
		Channel: "be04afda-4308-44bf-9a8d-633a93875270",
	}

	err := client.SendSentMessage(flowsMessage)
	if err != nil {
		t.Errorf("Error sending message: %s", err.Error())
	}
}
