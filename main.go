package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

var (
	bot *tgbotapi.BotAPI
)

func main() {
	log.Formatter = new(logrus.TextFormatter)
	log.Info("Sticker Cutter started!")

	var err error

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN env variable not specified!")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("authorized on account @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		// Handle only messages
		if update.Message == nil {
			continue
		}

		if update.Message.Photo != nil {
			HandleImage(update)
		}

		// commands
		switch update.Message.Command() {
		case "start":
			go StartCommand(update)
		}
	}
}
