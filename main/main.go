package main

import (
	"io"
	"log"
	"os"
	"telegram/youtube/bot/botconfig"
	"telegram/youtube/bot/validation"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	youtube "github.com/kkdai/youtube/v2"
)

func ExampleClient(link string) (string, string) {
	client := youtube.Client{}

	video, err := client.GetVideo(link)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[3])
	if err != nil {
		panic(err)
	}
	filepath := "/home/slovac/src/telegram-mp3-bot/audios/" + link + ".mp4"
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	return filepath, video.Title
}

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
				filePath, title := ExampleClient(link)
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
