package handler

import (
	"discord-proxy-flows/config"
	"discord-proxy-flows/internal/flows"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	Conf        config.Config
	DG          discordgoSessionInterface
	Server      *http.Server
	FlowsClient *flows.Client
}

func (h Handler) ReceiveMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	receivedMessage := flows.Message{
		From: m.Author.ID,
		Text: m.Content,
	}

	if len(m.Attachments) > 0 {
		receivedMessage.Attachments = []string{}
		for _, attachment := range m.Attachments {
			receivedMessage.Attachments = append(receivedMessage.Attachments, attachment.URL)
		}
	}

	if err := h.FlowsClient.SendReceiveMessage(receivedMessage); err != nil {
		log.Println("error on send message from discord bot to flows:", err)
	}
	log.Println("forwarded message to flows:", receivedMessage.From, receivedMessage.Text, receivedMessage.Attachments)
}

type SendMessageRequest struct {
	Message     string
	RecipientID string
}

func (h Handler) SetupHandler() {
	http.HandleFunc("/discord/rp/send", h.SendDirectMessage)
	h.Server.Handler = http.DefaultServeMux
}

func (h Handler) SendDirectMessage(w http.ResponseWriter, r *http.Request) {
	var msgReq flows.FlowsMessage
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Println("error on decode request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dmChannel, err := h.DG.UserChannelCreate(msgReq.To)
	if err != nil {
		log.Println("error on stablish dm channel connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, attachment := range msgReq.Attachments {
		req, err := http.NewRequest(http.MethodGet, attachment, nil)
		if err != nil {
			log.Println("Failed to get attachment")
		}

		filename := path.Base(attachment)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("error on get attachment file:", err)
			continue
		}

		_, err = h.DG.ChannelFileSend(dmChannel.ID, filename, res.Body)
		if err != nil {
			log.Println("error on send attachment to discord")
		}
		log.Println("sent attachment to discord bot", attachment)
	}

	if msgReq.Text == "" {
		return
	}
	_, err = h.DG.ChannelMessageSend(dmChannel.ID, msgReq.Text)
	if err != nil {
		log.Println("error on send message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("sent message to discord bot with ID:", msgReq.ID)

	h.FlowsClient.SendSentMessage(msgReq)

	w.WriteHeader(http.StatusOK)
}

type discordgoSessionInterface interface {
	UserChannelCreate(string, ...discordgo.RequestOption) (*discordgo.Channel, error)
	ChannelMessageSend(string, string, ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSend(string, string, io.Reader, ...discordgo.RequestOption) (*discordgo.Message, error)
}

type FlowsMessage struct {
	ID          string   `json:"id,omitempty"`
	Text        string   `json:"text,omitempty"`
	To          string   `json:"to,omitempty"`
	Channel     string   `json:"channel,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}
