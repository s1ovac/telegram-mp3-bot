package main

import (
	"log"
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
			link, err := url.ParseRequestURI(update.Message.Text)
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, link.Host)
			bot.Send(msg)
		}
	}
}
