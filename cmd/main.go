package main

import (
	"discord-proxy-flows/config"
	"discord-proxy-flows/internal/flows"
	"discord-proxy-flows/internal/handler"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error on load configuration:", err)
		log.Fatal(err)
		return
	}

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		fmt.Println("Error on create discord session:", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error on open discord connection:", err)
		return
	}
	defer dg.Close()

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	client := flows.NewClient(cfg.FlowsURL + "/c/ds/" + cfg.ChannelUUID)

	h := &handler.Handler{
		Conf:        cfg,
		DG:          dg,
		Server:      server,
		FlowsClient: client,
	}

	dg.AddHandler(h.ReceiveMessage)

	h.SetupHandler()

	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		fmt.Println("Error on start http server:", err)
		return
	}
}
