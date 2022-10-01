package youtubeapi

import (
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
)

func SaveAudio(link string) (string, string) {
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
