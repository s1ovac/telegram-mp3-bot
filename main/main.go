package main

import (
	"log"
	"net/http"
	"net/url"
	"telegram/bot/mp3/botconfig"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(botconfig.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			_, err := url.ParseRequestURI(update.Message.Text)
			if err != nil {
				log.Fatal(err)
			}
			response, err := http.Head(update.Message.Text)
			if err != nil {
				errorMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your link is not valid!\n Try more...")
				bot.Send(errorMsg)
				log.Fatal(err)
			}
			if response.Status != ("200 OK") {
				errorMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "I can't find video by this link\nTry more...")
				bot.Send(errorMsg)
				log.Fatal(response.Status)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.Status)
			bot.Send(msg)
		}
	}
}
