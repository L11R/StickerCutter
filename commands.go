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
	"net"
	"golang.org/x/net/proxy"
	"fmt"
	"context"
	"github.com/spf13/viper"
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

			res, err := http.Get(url)
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
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				socksDialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", viper.GetString("proxy.address"), viper.GetString("proxy.port")), &proxy.Auth{
					User: viper.GetString("proxy.user"),
					Password: viper.GetString("proxy.password"),
				}, proxy.Direct)
				if err != nil {
					return nil, err
				}

				return socksDialer.Dial(network, addr)
			},
		},
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