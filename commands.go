package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"image"
	"bytes"
	"image/png"
)

func StartCommand(update tgbotapi.Update)  {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me image.")
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

func HandleImage(update tgbotapi.Update) {
	for i, photo := range *update.Message.Photo {
		if i+1 == len(*update.Message.Photo) {
			url, err := bot.GetFileDirectURL(photo.FileID)
			if err != nil {
				log.Warn(err)
				return
			}

			res, err := http.Get(url)
			if err != nil {
				log.Warn(err)
				return
			}

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
