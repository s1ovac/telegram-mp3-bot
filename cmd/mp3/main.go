package main

import (
	"log"
	"os"
	"telegram/youtube/bot/internal/config"
	"telegram/youtube/bot/internal/validation"
	"telegram/youtube/bot/internal/youtubeapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую, "+update.Message.Chat.LastName+"\nОтправьте ссылку на видео...\n")
					bot.Send(msg)
				case "help":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для работы необходимо отправить ссылку на видео в формате «youtube.com/... или «youtu.be/...»\n")
					bot.Send(msg)
				case "commands":
					for command, description := range config.Commands {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, command+" - "+description)
						bot.Send(msg)
					}
				case "info":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Создатель: s1ovac\nGit: https://github.com/s1ovac\nВК: https://vk.com/slovacccc\n\nВсе права защищены.\nCopyright © 2022 «s1ovac»\n")
					bot.Send(msg)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда!\nСписок команд: /help")
					bot.Send(msg)
				}

			}
			if !update.Message.IsCommand() {
				err := validation.Validation(update.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильная ссылка!\nПопробуйте еще...")
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
					err = os.Remove(filePath)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}

	}
}
