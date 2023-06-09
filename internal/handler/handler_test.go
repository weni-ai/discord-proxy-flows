package handler

import (
	"bytes"
	"discord-proxy-flows/config"
	"discord-proxy-flows/internal/flows"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHandleSendMessage(t *testing.T) {

	dg := &discordgoMockSession{}

	requestBody, _ := json.Marshal(SendMessageRequest{
		Message:     "Hello",
		RecipientID: "112233445566778899",
	})
	request := httptest.NewRequest("POST", "/discord/rp/send", bytes.NewBuffer(requestBody))
	responseRecorder := httptest.NewRecorder()

	testServer := ServerTest(t)
	defer testServer.Close()

	client := flows.NewClient(testServer.URL)

	h := &Handler{
		Conf:        config.Config{},
		DG:          dg,
		Server:      nil,
		FlowsClient: client,
	}

	h.SendDirectMessage(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Esperado código de status %d, mas obteve %d", http.StatusOK, responseRecorder.Code)
	}
}

type discordgoMockSession struct{}

func (dgs *discordgoMockSession) UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (*discordgo.Channel, error) {
	if recipientID == "mocked_id_error" {
		return nil, errors.New("mocked error")
	}

	return &discordgo.Channel{ID: "mocked_id"}, nil
}

func (dgs *discordgoMockSession) ChannelMessageSend(dmChannelID, message string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	if dmChannelID == "mocked_dm_channel_id_error" {
		return nil, errors.New("mocked error")
	}

	return &discordgo.Message{}, nil
}

func (dgs *discordgoMockSession) ChannelFileSend(dmCHannelID, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return nil, nil
}

func ServerTest(t *testing.T) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sent" {
			t.Errorf("Expected route /sent, but got %s", r.URL.Path)
		}

		var receivedMessage flows.Message
		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		if err != nil {
			t.Errorf("Error decoding request: %s", err.Error())
		}

		if receivedMessage.Text != "Hello" {
			t.Errorf("Expected message 'Hello', but got %s", receivedMessage.Text)
		}

		if receivedMessage.From != "112233445566778899" {
			t.Errorf("Expected sender '112233445566778899', but got %s", receivedMessage.From)
		}

		w.WriteHeader(http.StatusOK)
	}))

	return testServer
}
