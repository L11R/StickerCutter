package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"image"
	"bytes"
	"image/png"
	"os"
	"io"
	"io/ioutil"
)

func StartCommand(update tgbotapi.Update)  {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me image.")
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func HandlePhoto(update tgbotapi.Update) {
	for i, photo := range *update.Message.Photo {
		if i+1 == len(*update.Message.Photo) {
			url, err := bot.GetFileDirectURL(photo.FileID)
			if err != nil {
				log.Warn(err)
				return
			}

			client := http.Client{
				Transport: tr,
			}

			res, err := client.Get(url)
			if err != nil {
				log.Warn(err)
				return
			}
			defer res.Body.Close()

			img, _, err := image.Decode(res.Body)
			if err != nil {
				log.Warn(err)
				return
			}
			croppedImg := MakeSticker(img)

			buf := new(bytes.Buffer)
			err = png.Encode(buf, croppedImg)
			if err != nil {
				log.Warn(err)
				return
			}

			msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
				Name: "sticker.png",
				Bytes: buf.Bytes(),
			})
			bot.Send(msg)
		}
	}
}

func HandleVideoNote(update tgbotapi.Update) {
	url, err := bot.GetFileDirectURL(update.Message.VideoNote.FileID)
	if err != nil {
		log.Warn(err)
		return
	}

	client := http.Client{
		Transport: tr,
	}

	res, err := client.Get(url)
	if err != nil {
		log.Warn(err)
		return
	}
	defer res.Body.Close()

	// create "temp1.mp4"
	file, err := os.Create("temp1.mp4")
	if err != nil {
		log.Warn(err)
		return
	}

	io.Copy(file, res.Body)

	err = RemoveAudio(res.Body)
	if err != nil {
		log.Warn(err)
		return
	}

	// open "temp2.mp4"
	fileBytes, err := ioutil.ReadFile("temp2.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name: "video.mp4",
		Bytes: fileBytes,
	})
	bot.Send(msg)
}