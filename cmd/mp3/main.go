package main

import (
	"log"
	"telegram/youtube/bot/internal/config"
	"telegram/youtube/bot/internal/config/validation"
	"telegram/youtube/bot/internal/config/youtubeapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
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
			err := validation.Validation(update.Message.Text)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Link is invalid!\nTry more...")
				bot.Send(msg)
			} else {
				link := update.Message.Text
				err = validation.Url(&link)
				if err != nil {
					log.Fatal(err)
				}
				filePath, title := youtubeapi.SaveAudio(link)
				file := tgbotapi.FilePath(filePath)
				if err != nil {
					log.Fatal(err)
				}
				msg := tgbotapi.NewAudio(update.Message.Chat.ID, file)
				msg.Title = title
				bot.Send(msg)
			}
		}
	}
}
